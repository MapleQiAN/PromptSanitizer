<template>
  <div style="display: flex; flex-direction: column; gap: 20px;">
    <!-- Processing Mode -->
    <div>
      <label>Mode</label>
      <select
        :value="config.mode"
        @change="(e) => updateConfig({ mode: (e.target as HTMLSelectElement).value as any })"
      >
        <option value="sanitize">Sanitize</option>
        <option value="annotate">Annotate</option>
      </select>
    </div>

    <!-- Strategy -->
    <div>
      <label>Strategy</label>
      <select
        :value="config.strategy"
        @change="(e) => updateConfig({ strategy: (e.target as HTMLSelectElement).value as any })"
      >
        <option value="redact">Redact (Placeholder)</option>
        <option value="mask">Mask (Partial)</option>
        <option value="pseudonym">Pseudonym (Consistent)</option>
      </select>
    </div>

    <!-- Detection Level -->
    <div>
      <label>Detection Level</label>
      <select
        :value="config.level"
        @change="(e) => updateConfig({ level: (e.target as HTMLSelectElement).value as any })"
      >
        <option value="lenient">Lenient</option>
        <option value="standard">Standard</option>
        <option value="strict">Strict</option>
      </select>
    </div>

    <!-- Category Toggles -->
    <div>
      <label>Detection Categories</label>
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
      <label>Allowlist (one per line)</label>
      <textarea
        :value="config.allowlist.join('\n')"
        @input="(e) => updateConfig({ allowlist: (e.target as HTMLTextAreaElement).value.split('\n').filter((s) => s.trim()) })"
        placeholder="Enter strings to exclude..."
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
      <label>Font Size (px)</label>
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
import type { Config } from "../types";

interface Props {
  config: Config;
  onChange: (config: Config) => void;
}

const props = defineProps<Props>();

const categories = [
  { id: "phone", label: "Phone" },
  { id: "email", label: "Email" },
  { id: "id_card", label: "ID Card" },
  { id: "ip", label: "IP Address" },
  { id: "domain", label: "Domain/URL" },
  { id: "token", label: "Token/Key" },
  { id: "password", label: "Password" },
  { id: "private_key", label: "Private Key" },
];

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
