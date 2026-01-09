<template>
  <div style="display: flex; flex-direction: column; height: 100vh">
    <!-- 顶部工具栏 -->
    <div
      style="
        padding: 12px 16px;
        background: #fff;
        border-bottom: 1px solid #e0e0e0;
        display: flex;
        gap: 8px;
        align-items: center;
      "
    >
      <button
        @click="handleSanitize"
        style="
          background: #4a90e2;
          color: white;
          font-weight: bold;
        "
      >
        一键清洗
      </button>
      <button @click="handleLoadFile">导入文件</button>
      <button @click="showConfig = !showConfig">
        {{ showConfig ? "隐藏配置" : "显示配置" }}
      </button>
      <div style="flex: 1" />
      <select
        :value="viewMode"
        @change="(e) => viewMode = (e.target as HTMLSelectElement).value as any"
        style="padding: 6px 12px"
      >
        <option value="split">分栏视图</option>
        <option value="diff">对比视图</option>
        <option value="report">报告视图</option>
      </select>
    </div>

    <!-- 配置面板 -->
    <ConfigPanel v-if="showConfig" :config="config" :onChange="setConfig" />

    <!-- 主内容区 -->
    <div style="flex: 1; display: flex; overflow: hidden">
      <template v-if="viewMode === 'split'">
        <MainPanel
          title="原文"
          :text="originalText"
          :onChange="setOriginalText"
          :findings="findings"
          :highlightMode="config.mode === 'annotate'"
        />
        <div style="width: 1px; background: #e0e0e0" />
        <MainPanel
          title="清洗结果"
          :text="sanitizedText"
          :readOnly="true"
          :findings="findings"
        />
      </template>

      <div v-if="viewMode === 'diff'" style="flex: 1; padding: 16px; overflow: auto">
        <h3 style="margin-bottom: 16px">前后对比</h3>
        <div
          style="
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 16px;
          "
        >
          <div>
            <h4>原文</h4>
            <pre
              style="
                background: #fff;
                padding: 12px;
                border: 1px solid #ddd;
                border-radius: 4px;
                white-space: pre-wrap;
                font-size: 13px;
              "
            >
              {{ originalText || "(空)" }}
            </pre>
          </div>
          <div>
            <h4>清洗后</h4>
            <pre
              style="
                background: #fff;
                padding: 12px;
                border: 1px solid #ddd;
                border-radius: 4px;
                white-space: pre-wrap;
                font-size: 13px;
              "
            >
              {{ sanitizedText || "(空)" }}
            </pre>
          </div>
        </div>
      </div>

      <ReportView v-if="viewMode === 'report' && response" :response="response" />
    </div>

    <!-- 底部命中列表 -->
    <FindingsList
      v-if="findings.length > 0"
      :findings="findings"
      :originalText="originalText"
      :onJump="(start, end) => {
        // 滚动到指定位置（简化实现）
        console.log('跳转到:', start, end);
      }"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { invoke } from "@tauri-apps/api/core";
import MainPanel from "./components/MainPanel.vue";
import ConfigPanel from "./components/ConfigPanel.vue";
import FindingsList from "./components/FindingsList.vue";
import ReportView from "./components/ReportView.vue";
import type { Config, Finding, Response } from "./types";

const originalText = ref("");
const sanitizedText = ref("");
const findings = ref<Finding[]>([]);
const response = ref<Response | null>(null);
const config = ref<Config>({
  mode: "sanitize",
  strategy: "redact",
  level: "standard",
  enabledCategories: [
    "phone",
    "email",
    "id_card",
    "ip",
    "domain",
    "token",
    "password",
    "private_key",
  ],
  allowlist: [],
});
const showConfig = ref(false);
const viewMode = ref<"split" | "diff" | "report">("split");

const setOriginalText = (text: string) => {
  originalText.value = text;
};

const setConfig = (newConfig: Config) => {
  config.value = newConfig;
};

const handleSanitize = async () => {
  if (!originalText.value.trim()) {
    alert("请输入要清洗的文本");
    return;
  }

  try {
    const result = await invoke<Response>("sanitize", {
      request: {
        text: originalText.value,
        mode: config.value.mode,
        strategy: config.value.strategy,
        level: config.value.level,
        enabled_categories: config.value.enabledCategories,
        allowlist: config.value.allowlist,
        semantic_mode: "off",
      },
    });

    response.value = result;
    sanitizedText.value = result.sanitized_text;
    findings.value = result.findings;
  } catch (error) {
    console.error("清洗失败:", error);
    alert(`清洗失败: ${error}`);
  }
};

const handleLoadFile = async () => {
  try {
    const filePath = await invoke<string>("open_file_dialog");
    if (filePath) {
      const content = await invoke<string>("read_file", { path: filePath });
      originalText.value = content;
    }
  } catch (error) {
    console.error("加载文件失败:", error);
  }
};
</script>
