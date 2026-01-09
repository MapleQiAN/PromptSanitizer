<template>
  <div style="display: flex; flex-direction: column; height: 100%;">
    <!-- Badge/Status -->
    <div v-if="highlightMode" style="margin-bottom: 16px;">
      <span class="badge">
        ⚠ Highlight Mode Active
      </span>
    </div>

    <!-- Highlighted View -->
    <div
      v-if="highlightMode && findings.length > 0"
      class="highlight-view"
      style="
        flex: 1;
        padding: 20px;
        background: var(--color-bg-tertiary);
        border: var(--border-thin) solid var(--color-border);
        overflow: auto;
      "
    >
      <span
        v-for="(part, i) in highlightedParts"
        :key="i"
        :class="{
          'highlight': part.isFinding,
          'highlight--high': part.isFinding && part.risk && part.risk >= 70,
          'highlight--low': part.isFinding && part.risk && part.risk < 40,
        }"
      >
        {{ part.text }}
      </span>
    </div>

    <!-- Text Editor -->
    <textarea
      v-else
      ref="textareaRef"
      :value="text"
      @input="handleInput"
      :readonly="readOnly"
      :placeholder="readOnly ? 'Sanitized output will appear here...' : 'Input or paste text for privacy sanitization...'"
      class="text-editor"
      :style="{ flex: 1 }"
    />

    <!-- Action Bar -->
    <div
      v-if="text"
      style="
        margin-top: 16px;
        padding-top: 16px;
        border-top: var(--border-thin) solid var(--color-border);
        display: flex;
        justify-content: flex-end;
        gap: 12px;
      "
    >
      <button
        class="btn-action"
        @click="handleCopy"
        style="padding: 8px 20px; font-size: 11px;"
      >
        📋 Copy {{ readOnly ? 'Output' : 'Text' }}
      </button>
      <div
        v-if="text"
        style="
          font-family: var(--font-mono);
          font-size: 10px;
          color: var(--color-text-muted);
          display: flex;
          align-items: center;
          padding: 0 8px;
        "
      >
        {{ text.length }} chars
      </div>
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
