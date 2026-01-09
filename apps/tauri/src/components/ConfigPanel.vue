<template>
  <div style="display: flex; flex-direction: column; gap: 20px;">
    <!-- Processing Mode -->
    <div>
      <label>{{ t.mode }}</label>
      <CustomSelect
        :model-value="config.mode"
        :options="[
          { value: 'sanitize', label: t.sanitize },
          { value: 'annotate', label: t.annotate },
        ]"
        @update:model-value="(value) => updateConfig({ mode: value as any })"
      />
    </div>

    <!-- Strategy -->
    <div>
      <label>{{ t.strategy }}</label>
      <CustomSelect
        :model-value="config.strategy"
        :options="[
          { value: 'redact', label: t.redact },
          { value: 'mask', label: t.mask },
          { value: 'pseudonym', label: t.pseudonym },
        ]"
        @update:model-value="(value) => updateConfig({ strategy: value as any })"
      />
    </div>

    <!-- Detection Level -->
    <div>
      <label>{{ t.level }}</label>
      <CustomSelect
        :model-value="config.level"
        :options="[
          { value: 'lenient', label: t.lenient },
          { value: 'standard', label: t.standard },
          { value: 'strict', label: t.strict },
        ]"
        @update:model-value="(value) => updateConfig({ level: value as any })"
      />
    </div>

    <!-- Category Toggles -->
    <div>
      <label>{{ t.categories }}</label>
      <div
        style="
          display: grid;
          grid-template-columns: repeat(2, 1fr);
          gap: 10px;
          margin-top: 8px;
        "
      >
        <label
          v-for="cat in categories"
          :key="cat.id"
          style="
            display: flex;
            align-items: center;
            gap: 10px;
            cursor: pointer;
            font-size: 12px;
            font-weight: 500;
            text-transform: none;
            letter-spacing: normal;
            color: var(--color-text-primary);
          "
        >
          <input
            type="checkbox"
            :checked="config.enabledCategories.includes(cat.id)"
            @change="toggleCategory(cat.id)"
          />
          <span>{{ cat.label }}</span>
        </label>
      </div>
    </div>

    <!-- Allowlist -->
    <div>
      <label>{{ t.allowlist }}</label>
      <textarea
        :value="config.allowlist.join('\n')"
        @input="(e) => updateConfig({ allowlist: (e.target as HTMLTextAreaElement).value.split('\n').filter((s) => s.trim()) })"
        :placeholder="t.allowlistPlaceholder"
        style="
          width: 100%;
          min-height: 80px;
          padding: 16px;
          background: var(--color-bg-secondary);
          border: 3px solid var(--color-border);
          border-radius: var(--radius-md);
          color: var(--color-text-primary);
          font-family: var(--font-mono);
          font-size: 13px;
          line-height: 1.6;
          resize: vertical;
          outline: none;
          transition: all 0.3s var(--ease-smooth);
        "
      />
    </div>

    <!-- Font Size -->
    <div>
      <label>{{ t.fontSize }}</label>
      <input
        type="number"
        :value="config.fontSize || 16"
        @input="(e) => updateConfig({ fontSize: parseInt((e.target as HTMLInputElement).value) || 16 })"
        min="10"
        max="32"
        step="1"
        style="
          width: 100%;
          padding: calc(var(--unit) * 2);
          background: var(--color-bg-secondary);
          border: 3px solid var(--color-border);
          outline: none;
          transition: all 0.3s var(--ease-smooth);
          border-radius: var(--radius-md);
          font-weight: 500;
        "
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, toRef } from "vue";
import CustomSelect from "./CustomSelect.vue";
import type { Config } from "../types";

type Language = "zh" | "en";

interface Props {
  config: Config;
  onChange: (config: Config) => void;
  lang?: Language;
}

const props = withDefaults(defineProps<Props>(), {
  lang: "zh",
});

// 使用 toRef 确保响应式追踪
const langRef = toRef(props, "lang");

const i18n = {
  zh: {
    mode: "处理模式",
    sanitize: "净化",
    annotate: "标注",
    strategy: "处理策略",
    redact: "遮蔽（占位符）",
    mask: "掩码（部分）",
    pseudonym: "假名（一致）",
    level: "检测级别",
    lenient: "宽松",
    standard: "标准",
    strict: "严格",
    categories: "检测类别",
    allowlist: "白名单（每行一个）",
    allowlistPlaceholder: "输入要排除的字符串...",
    fontSize: "字体大小 (px)",
    phone: "电话",
    email: "邮箱",
    id_card: "身份证",
    ip: "IP地址",
    domain: "域名/URL",
    token: "令牌/密钥",
    password: "密码",
    private_key: "私钥",
  },
  en: {
    mode: "Mode",
    sanitize: "Sanitize",
    annotate: "Annotate",
    strategy: "Strategy",
    redact: "Redact (Placeholder)",
    mask: "Mask (Partial)",
    pseudonym: "Pseudonym (Consistent)",
    level: "Detection Level",
    lenient: "Lenient",
    standard: "Standard",
    strict: "Strict",
    categories: "Detection Categories",
    allowlist: "Allowlist (one per line)",
    allowlistPlaceholder: "Enter strings to exclude...",
    fontSize: "Font Size (px)",
    phone: "Phone",
    email: "Email",
    id_card: "ID Card",
    ip: "IP Address",
    domain: "Domain/URL",
    token: "Token/Key",
    password: "Password",
    private_key: "Private Key",
  },
};

// 使用 langRef.value 确保响应式更新
const t = computed(() => i18n[langRef.value]);

const categories = computed(() => [
  { id: "phone", label: t.value.phone },
  { id: "email", label: t.value.email },
  { id: "id_card", label: t.value.id_card },
  { id: "ip", label: t.value.ip },
  { id: "domain", label: t.value.domain },
  { id: "token", label: t.value.token },
  { id: "password", label: t.value.password },
  { id: "private_key", label: t.value.private_key },
]);

const toggleCategory = (categoryId: string) => {
  const enabled = props.config.enabledCategories.includes(categoryId);
  props.onChange({
    ...props.config,
    enabledCategories: enabled
      ? props.config.enabledCategories.filter((c) => c !== categoryId)
      : [...props.config.enabledCategories, categoryId],
  });
};

const updateConfig = (updates: Partial<Config>) => {
  props.onChange({ ...props.config, ...updates });
};
</script>
