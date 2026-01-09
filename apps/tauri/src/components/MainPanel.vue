<template>
  <div style="flex: 1; display: flex; flex-direction: column; overflow: hidden">
    <div
      style="
        padding: 12px 16px;
        background: #f8f8f8;
        border-bottom: 1px solid #e0e0e0;
        font-weight: bold;
      "
    >
      {{ title }}
    </div>
    <div style="flex: 1; padding: 16px; overflow: auto">
      <div
        v-if="highlightMode && findings.length > 0"
        style="
          padding: 12px;
          background: #fff;
          border: 1px solid #ddd;
          border-radius: 4px;
          min-height: 200px;
          font-size: 13px;
          line-height: 1.6;
          white-space: pre-wrap;
        "
      >
        <span
          v-for="(part, i) in highlightedParts"
          :key="i"
          :style="{
            background:
              part.isFinding && part.risk
                ? part.risk >= 70
                  ? '#ffebee'
                  : part.risk >= 40
                  ? '#fff3e0'
                  : '#e3f2fd'
                : 'transparent',
            padding: part.isFinding ? '2px 4px' : '0',
            borderRadius: part.isFinding ? '2px' : '0',
          }"
        >
          {{ part.text }}
        </span>
      </div>
      <textarea
        v-else
        ref="textareaRef"
        :value="text"
        @input="handleInput"
        :readonly="readOnly"
        placeholder="在此输入或粘贴文本..."
        :style="{
          width: '100%',
          height: '100%',
          minHeight: '400px',
          padding: '12px',
          border: '1px solid #ddd',
          borderRadius: '4px',
          fontSize: '13px',
          fontFamily: 'monospace',
          resize: 'none',
          background: readOnly ? '#f9f9f9' : '#fff',
        }"
      />
    </div>
    <div v-if="!readOnly" style="padding: 8px 16px; border-top: 1px solid #e0e0e0">
      <button @click="handleCopy" :disabled="!text">复制</button>
    </div>
    <div v-if="readOnly && text" style="padding: 8px 16px; border-top: 1px solid #e0e0e0">
      <button @click="handleCopy">复制结果</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import type { Finding } from "../types";

interface Props {
  title: string;
  text: string;
  onChange?: (text: string) => void;
  readOnly?: boolean;
  findings?: Finding[];
  highlightMode?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  readOnly: false,
  findings: () => [],
  highlightMode: false,
});

const textareaRef = ref<HTMLTextAreaElement | null>(null);

const highlightedParts = computed(() => {
  if (!props.findings.length || !props.highlightMode) {
    return [];
  }

  const parts: Array<{ text: string; isFinding: boolean; risk?: number }> = [];
  let lastIndex = 0;

  // 按位置排序
  const sortedFindings = [...props.findings].sort((a, b) => a.start - b.start);

  sortedFindings.forEach((finding) => {
    if (finding.start > lastIndex) {
      parts.push({
        text: props.text.substring(lastIndex, finding.start),
        isFinding: false,
      });
    }
    parts.push({
      text: props.text.substring(finding.start, finding.end),
      isFinding: true,
      risk: finding.risk,
    });
    lastIndex = finding.end;
  });

  if (lastIndex < props.text.length) {
    parts.push({
      text: props.text.substring(lastIndex),
      isFinding: false,
    });
  }

  return parts;
});

const handleInput = (e: Event) => {
  const target = e.target as HTMLTextAreaElement;
  props.onChange?.(target.value);
};

const handleCopy = () => {
  navigator.clipboard.writeText(props.text);
};
</script>
