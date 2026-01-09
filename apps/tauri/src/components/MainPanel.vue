<template>
  <div style="flex: 1; display: flex; flex-direction: column; overflow: hidden">
    <div class="card__header">
      <div class="card__title">
        {{ title }}
      </div>
      <span v-if="highlightMode" class="badge-soft">
        高亮命中片段
      </span>
    </div>
    <div class="card__body" style="flex: 1; overflow: auto">
      <div
        v-if="highlightMode && findings.length > 0"
        style="
          padding: 12px;
          background: radial-gradient(
            circle at top left,
            rgba(15, 23, 42, 0.96),
            rgba(15, 23, 42, 0.92)
          );
          border-radius: 12px;
          border: 1px solid rgba(31, 41, 55, 0.95);
          min-height: 200px;
          font-size: 12px;
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
                  ? 'rgba(248, 113, 113, 0.16)'
                  : part.risk >= 40
                  ? 'rgba(251, 191, 36, 0.16)'
                  : 'rgba(56, 189, 248, 0.14)'
                : 'transparent',
            padding: part.isFinding ? '1px 4px' : '0',
            borderRadius: part.isFinding ? '999px' : '0',
            border: part.isFinding
              ? part.risk && part.risk >= 70
                ? '1px solid rgba(248, 113, 113, 0.4)'
                : part.risk && part.risk >= 40
                ? '1px solid rgba(251, 191, 36, 0.4)'
                : '1px solid rgba(56, 189, 248, 0.5)'
              : 'none',
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
          minHeight: '340px',
          padding: '12px',
          borderRadius: '12px',
          fontSize: '12px',
          fontFamily: 'SF Mono, Menlo, Monaco, Consolas, Liberation Mono, Courier New, monospace',
          resize: 'none',
          background: readOnly
            ? 'radial-gradient(circle at top left, rgba(15, 23, 42, 0.96), rgba(15, 23, 42, 0.9))'
            : 'radial-gradient(circle at top left, rgba(15, 23, 42, 0.96), rgba(15, 23, 42, 0.92))',
          border: '1px solid rgba(31, 41, 55, 0.95)',
          color: '#e5e7eb',
        }"
      />
    </div>
    <div
      v-if="!readOnly"
      style="
        padding: 8px 16px 10px;
        border-top: 1px solid rgba(31, 41, 55, 0.95);
        display: flex;
        justify-content: flex-end;
      "
    >
      <button class="btn-ghost" @click="handleCopy" :disabled="!text">
        复制
      </button>
    </div>
    <div
      v-if="readOnly && text"
      style="
        padding: 8px 16px 10px;
        border-top: 1px solid rgba(31, 41, 55, 0.95);
        display: flex;
        justify-content: flex-end;
      "
    >
      <button class="btn-ghost" @click="handleCopy">复制结果</button>
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
