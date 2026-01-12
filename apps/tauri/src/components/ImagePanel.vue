<template>
  <div style="display: flex; flex-direction: column; height: 100%; gap: 12px;">
    <!-- Upload and Control Section -->
    <div style="display: flex; gap: 8px; flex-wrap: wrap;">
      <input
        ref="fileInputRef"
        type="file"
        accept="image/*"
        @change="handleFileSelect"
        style="display: none;"
      />
      <button class="btn-action" @click="triggerFileSelect" style="width: auto; white-space: nowrap;">
        📁 {{ lang === 'zh' ? '选择图片' : 'Select Image' }}
      </button>
      <button
        v-if="imageFile"
        class="btn-action"
        @click="handleDetectText"
        :disabled="isProcessing"
        style="width: auto; white-space: nowrap;"
      >
        {{ isProcessing ? '⏳' : '🔍' }} {{ isProcessing ? (lang === 'zh' ? '检测中...' : 'Detecting...') : (lang === 'zh' ? '检测文本' : 'Detect Text') }}
      </button>
      <button
        v-if="imageFile && detectedFindings.length > 0"
        class="btn-action btn-action--primary"
        @click="handleApplyMask"
        :disabled="isProcessing"
        style="width: auto; white-space: nowrap;"
      >
        🎨 {{ lang === 'zh' ? '应用打码' : 'Apply Mask' }}
      </button>
      <button
        v-if="maskedImageUrl"
        class="btn-action"
        @click="handleDownload"
        style="width: auto; white-space: nowrap;"
      >
        💾 {{ lang === 'zh' ? '下载' : 'Download' }}
      </button>
    </div>

    <!-- Image Display Area -->
    <div style="flex: 1; display: flex; align-items: center; justify-content: center; padding: 20px; overflow: auto; background: var(--color-bg-tertiary); border: 3px solid var(--color-border); border-radius: var(--radius-md); position: relative; min-height: 0;">
      <div v-if="!imageUrl" class="empty-state" style="text-align: center; color: var(--color-text-muted);">
        <div style="font-size: 48px; opacity: 0.2; margin-bottom: 12px;">🖼️</div>
        <div style="font-family: var(--font-display); font-size: 14px; font-weight: 600;">
          {{ lang === 'zh' ? '请选择一张图片' : 'Please select an image' }}
        </div>
      </div>
      <div v-else style="position: relative; max-width: 100%; max-height: 100%;">
        <img
          ref="imageRef"
          :src="imageUrl"
          alt="Image"
          style="max-width: 100%; max-height: 100%; object-fit: contain; display: block;"
          @load="handleImageLoad"
        />
        <canvas
          ref="overlayCanvasRef"
          style="position: absolute; top: 0; left: 0; pointer-events: none;"
        ></canvas>
      </div>
    </div>

    <!-- Findings List -->
    <div v-if="detectedFindings.length > 0" style="max-height: 120px; overflow-y: auto; padding: 12px; background: var(--color-bg-secondary); border: 3px solid var(--color-border); border-radius: var(--radius-md);">
      <div style="font-family: var(--font-display); font-size: 12px; font-weight: 600; margin-bottom: 8px; color: var(--color-text-secondary);">
        {{ lang === 'zh' ? '检测到的敏感信息' : 'Detected Sensitive Info' }} ({{ detectedFindings.length }})
      </div>
      <div style="display: flex; flex-wrap: wrap; gap: 6px;">
        <div
          v-for="(finding, index) in detectedFindings"
          :key="index"
          :style="{
            padding: '6px 10px',
            background: getRiskColor(finding.risk),
            color: '#FFFFFF',
            borderRadius: 'var(--radius-sm)',
            fontSize: '11px',
            fontFamily: 'var(--font-mono)',
            cursor: 'pointer',
            transition: 'all 0.2s',
          }"
          @mouseenter="highlightFinding(index)"
          @mouseleave="clearHighlight"
        >
          {{ finding.text }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from "vue";
import { createWorker } from "tesseract.js";
import { useMessage } from "naive-ui";
import type { ImageFinding } from "../types";

const message = useMessage();

interface Props {
  lang?: "zh" | "en";
  readOnly?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  lang: "zh",
  readOnly: false,
});

const emit = defineEmits<{
  (e: "maskedImage", url: string): void;
}>();

const fileInputRef = ref<HTMLInputElement | null>(null);
const imageRef = ref<HTMLImageElement | null>(null);
const overlayCanvasRef = ref<HTMLCanvasElement | null>(null);

const imageFile = ref<File | null>(null);
const imageUrl = ref<string>("");
const maskedImageUrl = ref<string>("");
const isProcessing = ref(false);
const detectedFindings = ref<ImageFinding[]>([]);
const worker = ref<any>(null);
const imageScale = ref({ x: 1, y: 1 });

const triggerFileSelect = () => {
  fileInputRef.value?.click();
};

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (file && file.type.startsWith("image/")) {
    imageFile.value = file;
    if (imageUrl.value) {
      URL.revokeObjectURL(imageUrl.value);
    }
    imageUrl.value = URL.createObjectURL(file);
    maskedImageUrl.value = "";
    detectedFindings.value = [];
    clearHighlight();
  } else {
    message.warning(props.lang === "zh" ? "请选择有效的图片文件" : "Please select a valid image file");
  }
};

const handleImageLoad = () => {
  if (imageRef.value && overlayCanvasRef.value) {
    const img = imageRef.value;
    const canvas = overlayCanvasRef.value;
    
    // 等待图片完全加载
    if (img.complete && img.naturalWidth > 0) {
      canvas.width = img.offsetWidth;
      canvas.height = img.offsetHeight;
      
      // 计算缩放比例
      imageScale.value = {
        x: img.offsetWidth / img.naturalWidth,
        y: img.offsetHeight / img.naturalHeight,
      };
      
      // 如果有检测结果，重新绘制覆盖层
      if (detectedFindings.value.length > 0) {
        drawOverlay();
      }
    }
  }
};

const handleDetectText = async () => {
  if (!imageFile.value) {
    message.warning(props.lang === "zh" ? "请先选择图片" : "Please select an image first");
    return;
  }

  isProcessing.value = true;
  detectedFindings.value = [];

  try {
    // 创建或重用 worker
    if (!worker.value) {
      worker.value = await createWorker({
        logger: (m: any) => {
          if (m.status === "recognizing text") {
            console.log(`OCR Progress: ${Math.round(m.progress * 100)}%`);
          }
        },
      });
      const langCode = props.lang === "zh" ? "chi_sim+eng" : "eng";
      await worker.value.loadLanguage(langCode);
      await worker.value.initialize(langCode);
    }

    // 执行 OCR
    const { data } = await worker.value.recognize(imageFile.value);
    
    // 处理检测结果
    const findings: ImageFinding[] = [];
    const sensitivePatterns = [
      { pattern: /\b\d{11}\b/g, type: "phone", risk: 80 }, // 手机号
      { pattern: /\b\d{18}\b/g, type: "id_card", risk: 90 }, // 身份证号
      { pattern: /\b\d{17}[\dXx]\b/g, type: "id_card", risk: 90 }, // 身份证号（带X）
      { pattern: /\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b/g, type: "email", risk: 70 }, // 邮箱
      { pattern: /\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b/g, type: "bank_card", risk: 85 }, // 银行卡号
    ];

    // 遍历所有识别的单词
    if (data.words) {
      for (const word of data.words) {
        const text = word.text.trim();
        if (!text) continue;

        // 检查是否匹配敏感信息模式
        for (const { pattern, type, risk } of sensitivePatterns) {
          if (pattern.test(text)) {
            findings.push({
              type,
              text,
              bbox: {
                x0: word.bbox.x0,
                y0: word.bbox.y0,
                x1: word.bbox.x1,
                y1: word.bbox.y1,
                text,
                confidence: word.confidence,
              },
              confidence: word.confidence,
              risk,
            });
            break;
          }
        }
      }
    }

    detectedFindings.value = findings;
    
    if (findings.length === 0) {
      message.info(props.lang === "zh" ? "未检测到敏感信息" : "No sensitive information detected");
    } else {
      message.success(
        props.lang === "zh" 
          ? `检测到 ${findings.length} 处敏感信息` 
          : `Detected ${findings.length} sensitive information`
      );
      drawOverlay();
    }
  } catch (error) {
    console.error("OCR failed:", error);
    message.error(
      props.lang === "zh" 
        ? `文本检测失败: ${error}` 
        : `Text detection failed: ${error}`
    );
  } finally {
    isProcessing.value = false;
  }
};

const drawOverlay = () => {
  if (!overlayCanvasRef.value || !imageRef.value) return;
  
  const canvas = overlayCanvasRef.value;
  const ctx = canvas.getContext("2d");
  if (!ctx) return;

  ctx.clearRect(0, 0, canvas.width, canvas.height);

  // 绘制检测到的区域
  detectedFindings.value.forEach((finding) => {
    const bbox = finding.bbox;
    const x = bbox.x0 * imageScale.value.x;
    const y = bbox.y0 * imageScale.value.y;
    const width = (bbox.x1 - bbox.x0) * imageScale.value.x;
    const height = (bbox.y1 - bbox.y0) * imageScale.value.y;

    // 绘制半透明矩形
    ctx.fillStyle = `rgba(255, 0, 0, 0.3)`;
    ctx.fillRect(x, y, width, height);
    
    // 绘制边框
    ctx.strokeStyle = `rgba(255, 0, 0, 0.8)`;
    ctx.lineWidth = 2;
    ctx.strokeRect(x, y, width, height);
  });
};

const highlightFinding = (index: number) => {
  if (!overlayCanvasRef.value || !imageRef.value) return;
  
  const canvas = overlayCanvasRef.value;
  const ctx = canvas.getContext("2d");
  if (!ctx) return;

  // 重新绘制所有区域，高亮选中的
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  
  detectedFindings.value.forEach((finding, i) => {
    const bbox = finding.bbox;
    const x = bbox.x0 * imageScale.value.x;
    const y = bbox.y0 * imageScale.value.y;
    const width = (bbox.x1 - bbox.x0) * imageScale.value.x;
    const height = (bbox.y1 - bbox.y0) * imageScale.value.y;

    if (i === index) {
      // 高亮选中的
      ctx.fillStyle = `rgba(255, 165, 0, 0.5)`;
      ctx.fillRect(x, y, width, height);
      ctx.strokeStyle = `rgba(255, 165, 0, 1)`;
      ctx.lineWidth = 3;
    } else {
      // 其他区域
      ctx.fillStyle = `rgba(255, 0, 0, 0.2)`;
      ctx.fillRect(x, y, width, height);
      ctx.strokeStyle = `rgba(255, 0, 0, 0.6)`;
      ctx.lineWidth = 2;
    }
    ctx.strokeRect(x, y, width, height);
  });
};

const clearHighlight = () => {
  drawOverlay();
};

const handleApplyMask = () => {
  if (!imageFile.value || !imageRef.value || detectedFindings.value.length === 0) {
    message.warning(props.lang === "zh" ? "请先检测文本" : "Please detect text first");
    return;
  }

  try {
    const img = imageRef.value;
    const canvas = document.createElement("canvas");
    canvas.width = img.naturalWidth;
    canvas.height = img.naturalHeight;
    const ctx = canvas.getContext("2d");
    
    if (!ctx) {
      throw new Error("Failed to get canvas context");
    }

    // 绘制原始图片
    ctx.drawImage(img, 0, 0);

    // 对检测到的区域进行打码（模糊处理）
    detectedFindings.value.forEach((finding) => {
      const bbox = finding.bbox;
      const x = Math.floor(bbox.x0);
      const y = Math.floor(bbox.y0);
      const width = Math.ceil(bbox.x1 - bbox.x0);
      const height = Math.ceil(bbox.y1 - bbox.y0);

      // 提取区域
      const imageData = ctx.getImageData(x, y, width, height);
      
      // 应用模糊效果（简单的像素化）
      const pixelSize = 8;
      for (let py = 0; py < height; py += pixelSize) {
        for (let px = 0; px < width; px += pixelSize) {
          let r = 0, g = 0, b = 0, a = 0, count = 0;
          
          for (let dy = 0; dy < pixelSize && py + dy < height; dy++) {
            for (let dx = 0; dx < pixelSize && px + dx < width; dx++) {
              const idx = ((py + dy) * width + (px + dx)) * 4;
              r += imageData.data[idx];
              g += imageData.data[idx + 1];
              b += imageData.data[idx + 2];
              a += imageData.data[idx + 3];
              count++;
            }
          }
          
          r = Math.floor(r / count);
          g = Math.floor(g / count);
          b = Math.floor(b / count);
          a = Math.floor(a / count);
          
          for (let dy = 0; dy < pixelSize && py + dy < height; dy++) {
            for (let dx = 0; dx < pixelSize && px + dx < width; dx++) {
              const idx = ((py + dy) * width + (px + dx)) * 4;
              imageData.data[idx] = r;
              imageData.data[idx + 1] = g;
              imageData.data[idx + 2] = b;
              imageData.data[idx + 3] = a;
            }
          }
        }
      }
      
      // 将处理后的图像数据放回
      ctx.putImageData(imageData, x, y);
    });

    // 转换为图片 URL
    maskedImageUrl.value = canvas.toDataURL("image/png");
    emit("maskedImage", maskedImageUrl.value);
    message.success(props.lang === "zh" ? "打码完成" : "Masking completed");
  } catch (error) {
    console.error("Masking failed:", error);
    message.error(
      props.lang === "zh" 
        ? `打码失败: ${error}` 
        : `Masking failed: ${error}`
    );
  }
};

const handleDownload = () => {
  if (!maskedImageUrl.value) {
    message.warning(props.lang === "zh" ? "没有可下载的图片" : "No image to download");
    return;
  }

  const link = document.createElement("a");
  link.href = maskedImageUrl.value;
  link.download = `masked_${imageFile.value?.name || "image.png"}`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  message.success(props.lang === "zh" ? "下载成功" : "Download successful");
};

const getRiskColor = (risk: number) => {
  if (risk >= 80) return "var(--color-risk-high, #ff4444)";
  if (risk >= 60) return "var(--color-risk-medium, #ffaa00)";
  return "var(--color-risk-low, #ffaa00)";
};

onUnmounted(() => {
  // 清理 worker
  if (worker.value) {
    worker.value.terminate();
  }
  // 清理图片 URL
  if (imageUrl.value) {
    URL.revokeObjectURL(imageUrl.value);
  }
});
</script>