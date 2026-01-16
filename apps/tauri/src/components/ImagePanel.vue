<template>
  <div style="display: flex; flex-direction: column; height: 100%; gap: 16px;">
    <!-- Upload and Control Section -->
    <div style="display: flex; flex-direction: column; gap: 12px;">
      <input
        ref="fileInputRef"
        type="file"
        accept="image/*"
        @change="handleFileSelect"
        style="display: none;"
      />
      <!-- Action Buttons Row -->
      <div style="display: flex; gap: 10px; align-items: center; flex-wrap: wrap;">
        <button
          v-if="imageFile"
          class="btn-action"
          @click="handleDetectText"
          :disabled="isProcessing"
          style="width: auto; white-space: nowrap; padding: 8px 16px;"
        >
          {{ isProcessing ? '⏳' : '🔍' }} {{ isProcessing ? (lang === 'zh' ? '检测中...' : 'Detecting...') : (lang === 'zh' ? '检测文本' : 'Detect Text') }}
        </button>
        <button
          v-if="imageFile && detectedFindings.length > 0"
          class="btn-action btn-action--primary"
          @click="handleApplyMask"
          :disabled="isProcessing"
          style="width: auto; white-space: nowrap; padding: 8px 20px;"
        >
          🎨 {{ lang === 'zh' ? '应用打码' : 'Apply Mask' }}
        </button>
        <button
          v-if="imageFile"
          class="btn-action"
          @click="handleClearImage"
          style="width: auto; white-space: nowrap; padding: 8px 16px;"
        >
          🗑️ {{ lang === 'zh' ? '清空图片' : 'Clear Image' }}
        </button>
        <button
          v-if="maskedImageUrl"
          class="btn-action"
          @click="handleDownload"
          style="width: auto; white-space: nowrap; padding: 8px 16px; margin-left: auto;"
        >
          💾 {{ lang === 'zh' ? '下载图片' : 'Download Image' }}
        </button>
      </div>
    </div>

    <!-- Image Display Area -->
    <div
      class="image-display-area"
      :class="{ 
        'image-display-area--dragging': isDragging,
        'image-display-area--clickable': !imageUrl
      }"
      @dragover.prevent.stop="handleDragOver"
      @dragenter.prevent.stop="handleDragOver"
      @dragleave.prevent="handleDragLeave"
      @drop.prevent.stop="handleDrop"
      @click="handleImageAreaClick"
    >
      <div v-if="!imageUrl" class="empty-state">
        <div class="empty-state__icon">🖼️</div>
        <div class="empty-state__title">
          {{ isDragging ? (lang === 'zh' ? '松开鼠标以加载图片' : 'Release to load image') : (lang === 'zh' ? '请选择一张图片或拖拽图片到这里' : 'Please select an image or drag and drop here') }}
        </div>
        <div class="empty-state__desc">
          {{ lang === 'zh' ? '支持 JPG、PNG 等常见图片格式' : 'Supports JPG, PNG and other common image formats' }}
        </div>
      </div>
      <div v-else style="position: relative; max-width: 100%; max-height: 100%; display: flex; align-items: center; justify-content: center;" @click.stop>
        <img
          ref="imageRef"
          :src="imageUrl"
          alt="Image"
          draggable="false"
          style="max-width: 100%; max-height: 100%; object-fit: contain; display: block; border-radius: var(--radius-sm);"
          @load="handleImageLoad"
          @dragover.prevent.stop="handleDragOver"
          @dragenter.prevent.stop="handleDragOver"
          @dragleave.prevent="handleDragLeave"
          @drop.prevent.stop="handleDrop"
          @click.stop
        />
        <canvas
          ref="overlayCanvasRef"
          style="position: absolute; top: 0; left: 0; pointer-events: none; border-radius: var(--radius-sm);"
        ></canvas>
      </div>
    </div>

    <!-- Findings List -->
    <div v-if="detectedFindings.length > 0" style="max-height: 140px; overflow-y: auto; padding: 16px; background: var(--color-bg-secondary); border: 3px solid var(--color-border); border-radius: var(--radius-md);">
      <div style="font-family: var(--font-display); font-size: 13px; font-weight: 600; margin-bottom: 12px; color: var(--color-text-secondary); display: flex; align-items: center; gap: 8px;">
        <span>🔍</span>
        <span>{{ lang === 'zh' ? '检测到的敏感信息' : 'Detected Sensitive Info' }}</span>
        <span style="background: var(--color-primary); color: #FFFFFF; padding: 2px 8px; border-radius: var(--radius-sm); font-size: 11px; font-weight: 700;">
          {{ detectedFindings.length }}
        </span>
      </div>
      <div style="display: flex; flex-wrap: wrap; gap: 8px;">
        <div
          v-for="(finding, index) in detectedFindings"
          :key="index"
          :style="{
            padding: '8px 12px',
            background: getRiskColor(finding.risk),
            color: '#FFFFFF',
            borderRadius: 'var(--radius-sm)',
            fontSize: '12px',
            fontFamily: 'var(--font-mono)',
            cursor: 'pointer',
            transition: 'all 0.2s',
            boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
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

<style scoped>
.image-display-area {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  overflow: auto;
  background: var(--color-bg-tertiary);
  border: 3px solid var(--color-border);
  border-radius: var(--radius-md);
  position: relative;
  min-height: 300px;
  transition: all 0.3s ease;
}

.image-display-area--dragging {
  border-color: var(--color-primary);
  background-color: var(--color-bg-secondary);
  border-width: 4px;
}

.image-display-area--clickable {
  cursor: pointer;
  transition: all 0.3s ease;
}

.image-display-area--clickable:hover {
  border-color: var(--color-primary);
  background-color: var(--color-bg-secondary);
}

.empty-state {
  text-align: center;
  color: var(--color-text-muted);
  padding: 40px 20px;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.empty-state__icon {
  font-size: 64px;
  opacity: 0.15;
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.image-display-area--dragging .empty-state__icon {
  transform: scale(1.1);
  opacity: 0.3;
}

.empty-state__title {
  font-family: var(--font-display);
  font-size: 15px;
  font-weight: 600;
}

.empty-state__desc {
  font-size: 13px;
  opacity: 0.7;
}
</style>

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
const isDragging = ref(false);

// 辅助函数 - 需要先定义
const getRiskColor = (risk: number) => {
  if (risk >= 80) return "var(--color-risk-high, #ff4444)";
  if (risk >= 60) return "var(--color-risk-medium, #ffaa00)";
  return "var(--color-risk-low, #ffaa00)";
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

// 事件处理函数
const triggerFileSelect = () => {
  fileInputRef.value?.click();
};

const handleImageAreaClick = () => {
  // 只在没有图片时允许点击选择文件
  // 避免在拖拽操作后立即触发点击
  if (!imageUrl.value && !isDragging.value) {
    // 检查是否是真正的点击（不是拖拽后的释放）
    const timeSinceDragEnd = Date.now() - ((window as any).lastDragEndTime || 0);
    if (!(window as any).lastDragEndTime || timeSinceDragEnd > 200) {
      triggerFileSelect();
    }
  }
};

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  processImageFile(file);
};

const processImageFile = (file: File | null | undefined) => {
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

const handleDragOver = (event: DragEvent) => {
  event.preventDefault();
  event.stopPropagation();
  // 检查是否是图片文件
  if (event.dataTransfer?.types.includes('Files')) {
    if (!isDragging.value) {
      isDragging.value = true;
    }
  }
};

const handleDragLeave = (event: DragEvent) => {
  event.preventDefault();
  // 只有当离开整个拖拽区域时才取消高亮
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect();
  const x = event.clientX;
  const y = event.clientY;
  
  if (x < rect.left || x > rect.right || y < rect.top || y > rect.bottom) {
    isDragging.value = false;
  }
};

const handleDrop = (event: DragEvent) => {
  event.preventDefault();
  event.stopPropagation();
  isDragging.value = false;
  
  // 记录拖拽结束时间，避免立即触发点击事件
  (window as any).lastDragEndTime = Date.now();

  const files = event.dataTransfer?.files;
  if (files && files.length > 0) {
    const file = files[0];
    console.log("拖拽文件:", file.name, file.type); // 调试信息
    if (file.type.startsWith("image/")) {
      processImageFile(file);
      message.success(props.lang === "zh" ? "图片加载成功" : "Image loaded successfully");
    } else {
      message.warning(props.lang === "zh" ? "请拖入有效的图片文件" : "Please drag a valid image file");
    }
  } else {
    console.log("没有检测到文件"); // 调试信息
    message.warning(props.lang === "zh" ? "未检测到文件" : "No file detected");
  }
};

const handleClearImage = () => {
  // 清理图片 URL
  if (imageUrl.value) {
    URL.revokeObjectURL(imageUrl.value);
    imageUrl.value = "";
  }
  // 重置所有状态
  imageFile.value = null;
  maskedImageUrl.value = "";
  detectedFindings.value = [];
  // 清空 canvas
  if (overlayCanvasRef.value) {
    const canvas = overlayCanvasRef.value;
    const ctx = canvas.getContext("2d");
    if (ctx) {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
    }
  }
  // 清空文件输入
  if (fileInputRef.value) {
    fileInputRef.value.value = "";
  }
  message.success(props.lang === "zh" ? "已清空图片" : "Image cleared");
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
      const langCode = props.lang === "zh" ? "chi_sim+eng" : "eng";
      worker.value = await createWorker(langCode);
    }

    // 执行 OCR
    const { data } = await worker.value.recognize(imageFile.value);
    
    // 处理检测结果
    const findings: ImageFinding[] = [];
    
    // 获取完整文本（用于地址等跨单词的检测）
    const fullText = data.text || "";
    
    // 改进的敏感信息检测模式
    const sensitivePatterns = [
      // 电话号码 - 支持多种格式
      { 
        pattern: /(\+?86[\s-]?)?1[3-9]\d[\s-]?\d{4}[\s-]?\d{4}/g, 
        type: "phone", 
        risk: 80,
        description: "手机号"
      },
      // 带括号区号的固定电话格式，如 (029) 8860 5585 或 (029) 88605585
      { 
        pattern: /\(0\d{2,3}\)[\s-]?\d{4}[\s-]?\d{4}/g, 
        type: "phone", 
        risk: 75,
        description: "固定电话（带括号区号）"
      },
      { 
        pattern: /(\+?86[\s-]?)?0\d{2,3}[\s-]?\d{7,8}/g, 
        type: "phone", 
        risk: 75,
        description: "固定电话"
      },
      { 
        pattern: /(\+?\d{1,3}[\s-]?)?\(?\d{3}\)?[\s-]?\d{3}[\s-]?\d{4}/g, 
        type: "phone", 
        risk: 75,
        description: "国际电话"
      },
      // 身份证号
      { 
        pattern: /\b\d{18}\b/g, 
        type: "id_card", 
        risk: 90,
        description: "身份证号"
      },
      { 
        pattern: /\b\d{17}[\dXx]\b/g, 
        type: "id_card", 
        risk: 90,
        description: "身份证号（带X）"
      },
      // 邮箱
      { 
        pattern: /\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b/g, 
        type: "email", 
        risk: 70,
        description: "邮箱"
      },
      // 银行卡号
      { 
        pattern: /\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}(\d{4})?/g, 
        type: "bank_card", 
        risk: 85,
        description: "银行卡号"
      },
      // 中文地址关键词模式 - 更灵活的匹配
      { 
        pattern: /[\u4e00-\u9fa5]{1,}(省|市|区|县|街道|路|街|巷|弄|号|小区|村|镇|乡|组|单元|室|层|栋|座|区|县)[\u4e00-\u9fa5\d\s-]*/g, 
        type: "address", 
        risk: 65,
        description: "中文地址"
      },
      // 更完整的中文地址（包含省市区）
      { 
        pattern: /[\u4e00-\u9fa5]{2,}省[\u4e00-\u9fa5]{1,}市[\u4e00-\u9fa5]{1,}(区|县)[\u4e00-\u9fa5\d\s-]*/g, 
        type: "address", 
        risk: 70,
        description: "完整中文地址"
      },
      // 省份+城市格式，如"陕西 西安/咸阳"、"山西 晋中"
      { 
        pattern: /[\u4e00-\u9fa5]{2,}[\s]+[\u4e00-\u9fa5]{1,}(?:\/[\u4e00-\u9fa5]{1,})?/g, 
        type: "address", 
        risk: 65,
        description: "省份城市位置"
      },
      // 英文地址模式
      { 
        pattern: /\d+\s+[A-Za-z\s]+(Street|St|Avenue|Ave|Road|Rd|Boulevard|Blvd|Lane|Ln|Drive|Dr|Court|Ct|Place|Pl|Way|Circle|Cir)[\s,]*[A-Z]{2}\s+\d{5}(-\d{4})?/gi, 
        type: "address", 
        risk: 65,
        description: "英文地址"
      },
      // 邮政编码（中国）- 更严格的匹配，避免误报
      { 
        pattern: /(邮编|邮政编码|邮编：|邮政编码：)[\s:：]*\d{6}/g, 
        type: "address", 
        risk: 50,
        description: "邮政编码"
      },
    ];

    // 首先在完整文本上检测跨单词的敏感信息（如地址、带分隔符的电话号）
    const textFindings: Array<{match: RegExpMatchArray, type: string, risk: number, description: string}> = [];
    for (const { pattern, type, risk, description } of sensitivePatterns) {
      const matches = fullText.matchAll(pattern);
      for (const match of matches) {
        if (match.index !== undefined) {
          textFindings.push({
            match,
            type,
            risk,
            description
          });
        }
      }
    }

    // 为每个匹配找到对应的单词边界框
    if (data.words && data.words.length > 0) {
      // 处理跨单词的匹配（如地址、带分隔符的电话号）
      for (const { match, type, risk } of textFindings) {
        const matchedText = match[0].trim();
        const matchStart = match.index!;
        const matchEnd = matchStart + match[0].length;
        
        // 方法1：通过文本位置匹配（更准确）
        // 重建文本并找到匹配的单词
        let textPos = 0;
        const relevantWords: typeof data.words = [];
        const matchedTextNormalized = matchedText.replace(/\s+/g, '').toLowerCase();
        
        for (const word of data.words) {
          const wordText = word.text;
          const wordStartInText = textPos;
          const wordEndInText = textPos + wordText.length;
          
          // 检查单词是否在匹配范围内
          if (wordEndInText > matchStart && wordStartInText < matchEnd) {
            relevantWords.push(word);
          }
          
          // 更新文本位置（考虑空格）
          textPos = wordEndInText;
          // 检查下一个字符是否是空格
          if (fullText[textPos] === ' ' || fullText[textPos] === '\n') {
            textPos++;
          }
        }
        
        // 方法2：如果方法1没找到，尝试通过文本内容匹配
        if (relevantWords.length === 0) {
          // 尝试找到包含匹配文本的单词
          for (const word of data.words) {
            const wordTextNormalized = word.text.replace(/\s+/g, '').toLowerCase();
            if (wordTextNormalized.includes(matchedTextNormalized) || 
                matchedTextNormalized.includes(wordTextNormalized)) {
              relevantWords.push(word);
            }
          }
        }
        
        if (relevantWords.length > 0) {
          // 计算合并的边界框
          const x0 = Math.min(...relevantWords.map((w: any) => w.bbox.x0));
          const y0 = Math.min(...relevantWords.map((w: any) => w.bbox.y0));
          const x1 = Math.max(...relevantWords.map((w: any) => w.bbox.x1));
          const y1 = Math.max(...relevantWords.map((w: any) => w.bbox.y1));
          const avgConfidence = relevantWords.reduce((sum: number, w: any) => sum + w.confidence, 0) / relevantWords.length;
          
          // 检查是否已经添加过（避免重复）
          const alreadyAdded = findings.some(f => 
            f.type === type &&
            f.text === matchedText &&
            Math.abs(f.bbox.x0 - x0) < 10 &&
            Math.abs(f.bbox.y0 - y0) < 10
          );
          
          if (!alreadyAdded) {
            findings.push({
              type,
              text: matchedText,
              bbox: {
                x0,
                y0,
                x1,
                y1,
                text: matchedText,
                confidence: avgConfidence,
              },
              confidence: avgConfidence,
              risk,
            });
          }
        }
      }

      // 也检查单个单词（用于简单的模式，如纯数字身份证号）
      for (const word of data.words) {
        const text = word.text.trim();
        if (!text) continue;

        // 只检查简单的单单词模式，避免重复检测
        const singleWordPatterns = [
          { pattern: /^\d{11}$/, type: "phone", risk: 80 },
          { pattern: /^\d{18}$/, type: "id_card", risk: 90 },
          { pattern: /^\d{17}[\dXx]$/, type: "id_card", risk: 90 },
          { pattern: /^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$/, type: "email", risk: 70 },
        ];

        for (const { pattern, type, risk } of singleWordPatterns) {
          if (pattern.test(text)) {
            // 检查是否已经添加过（避免重复）
            const alreadyAdded = findings.some(f => 
              f.text === text && 
              Math.abs(f.bbox.x0 - word.bbox.x0) < 5 &&
              Math.abs(f.bbox.y0 - word.bbox.y0) < 5
            );
            
            if (!alreadyAdded) {
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
            }
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

// 暴露方法给父组件
defineExpose({
  processImageFile,
});
</script>