<template>
  <div
    style="
      padding: 16px;
      background: #fff;
      border-bottom: 1px solid #e0e0e0;
    "
  >
    <h3 style="margin-bottom: 16px">配置选项</h3>
    <div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 16px">
      <!-- 模式 -->
      <div>
        <label style="display: block; margin-bottom: 8px; font-weight: bold">
          处理模式
        </label>
        <select
          :value="config.mode"
          @change="(e) => updateConfig({ mode: (e.target as HTMLSelectElement).value as any })"
          style="width: 100%; padding: 8px"
        >
          <option value="sanitize">清洗模式</option>
          <option value="annotate">标注模式</option>
        </select>
      </div>

      <!-- 清洗策略 -->
      <div>
        <label style="display: block; margin-bottom: 8px; font-weight: bold">
          清洗策略
        </label>
        <select
          :value="config.strategy"
          @change="(e) => updateConfig({ strategy: (e.target as HTMLSelectElement).value as any })"
          style="width: 100%; padding: 8px"
        >
          <option value="redact">占位符替换</option>
          <option value="mask">部分打码</option>
          <option value="pseudonym">一致化替换</option>
        </select>
      </div>

      <!-- 清洗强度 -->
      <div>
        <label style="display: block; margin-bottom: 8px; font-weight: bold">
          清洗强度
        </label>
        <select
          :value="config.level"
          @change="(e) => updateConfig({ level: (e.target as HTMLSelectElement).value as any })"
          style="width: 100%; padding: 8px"
        >
          <option value="lenient">宽松</option>
          <option value="standard">标准</option>
          <option value="strict">严格</option>
        </select>
      </div>

      <!-- 类别开关 -->
      <div>
        <label style="display: block; margin-bottom: 8px; font-weight: bold">
          检测类别
        </label>
        <div
          style="
            display: grid;
            grid-template-columns: repeat(2, 1fr);
            gap: 8px;
          "
        >
          <label
            v-for="cat in categories"
            :key="cat.id"
            style="
              display: flex;
              align-items: center;
              gap: 8px;
              cursor: pointer;
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
    </div>

    <!-- 白名单 -->
    <div style="margin-top: 16px">
      <label style="display: block; margin-bottom: 8px; font-weight: bold">
        白名单（每行一个，精确匹配）
      </label>
      <textarea
        :value="config.allowlist.join('\n')"
        @input="(e) => updateConfig({ allowlist: (e.target as HTMLTextAreaElement).value.split('\n').filter((s) => s.trim()) })"
        placeholder="输入要排除的字符串，每行一个"
        style="
          width: 100%;
          min-height: 80px;
          padding: 8px;
          border: 1px solid #ddd;
          border-radius: 4px;
          font-size: 13px;
        "
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import type { Config } from "../types";

interface Props {
  config: Config;
  onChange: (config: Config) => void;
}

const props = defineProps<Props>();

const categories = [
  { id: "phone", label: "手机号" },
  { id: "email", label: "邮箱" },
  { id: "id_card", label: "身份证" },
  { id: "ip", label: "IP地址" },
  { id: "domain", label: "域名/URL" },
  { id: "token", label: "Token/Key" },
  { id: "password", label: "密码" },
  { id: "private_key", label: "私钥" },
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
