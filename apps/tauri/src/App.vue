<template>
  <div class="app-shell">
    <!-- 顶部工具栏 / 品牌栏 -->
    <header class="app-shell__header">
      <div class="app-shell__brand">
        <div class="app-shell__avatar">
          <div class="app-shell__avatar-inner">PS</div>
        </div>
        <div class="app-shell__title">
          <span class="app-shell__title-main">PromptSanitizer</span>
          <span class="app-shell__title-sub">专注于智能脱敏与合规检查</span>
        </div>
      </div>

      <div class="app-shell__header-divider" />

      <div class="app-shell__toolbar">
        <button class="btn-primary" @click="handleSanitize">
          一键清洗
        </button>
        <button class="btn-ghost" @click="handleLoadFile">导入文件</button>
        <button class="btn-ghost" @click="showConfig = !showConfig">
          {{ showConfig ? "隐藏配置" : "显示配置" }}
        </button>

        <div style="flex: 1" />

        <select
          class="select-pill"
          :value="viewMode"
          @change="(e) => (viewMode = (e.target as HTMLSelectElement).value as any)"
        >
          <option value="split">分栏视图</option>
          <option value="diff">对比视图</option>
          <option value="report">报告视图</option>
        </select>
      </div>
    </header>

    <main class="app-shell__body">
      <!-- 配置面板 -->
      <section v-if="showConfig" class="card card--subtle">
        <div class="card__body">
          <ConfigPanel :config="config" :onChange="setConfig" />
        </div>
      </section>

      <!-- 主内容区 -->
      <section class="app-shell__main">
        <template v-if="viewMode === 'split'">
          <div style="flex: 1" class="card card--subtle">
            <MainPanel
              title="原文"
              :text="originalText"
              :onChange="setOriginalText"
              :findings="findings"
              :highlightMode="config.mode === 'annotate'"
            />
          </div>
          <div style="width: 1px; background: rgba(30, 64, 175, 0.4)" />
          <div style="flex: 1" class="card card--subtle">
            <MainPanel
              title="清洗结果"
              :text="sanitizedText"
              :readOnly="true"
              :findings="findings"
            />
          </div>
        </template>

        <section
          v-if="viewMode === 'diff'"
          class="card card--subtle"
          style="flex: 1; overflow: hidden"
        >
          <div class="card__header">
            <div class="card__title">前后对比</div>
            <span class="badge-soft">查看脱敏前后差异</span>
          </div>
          <div
            class="card__body"
            style="
              display: grid;
              grid-template-columns: 1fr 1fr;
              gap: 14px;
              height: 100%;
            "
          >
            <div style="display: flex; flex-direction: column; min-width: 0">
              <div
                style="
                  font-size: 12px;
                  text-transform: uppercase;
                  letter-spacing: 0.06em;
                  color: #9ca3af;
                  margin-bottom: 6px;
                "
              >
                原文
              </div>
              <pre
                style="
                  flex: 1;
                  background: radial-gradient(
                    circle at top left,
                    rgba(15, 23, 42, 0.96),
                    rgba(15, 23, 42, 0.92)
                  );
                  padding: 12px;
                  border-radius: 12px;
                  border: 1px solid rgba(31, 41, 55, 0.95);
                  white-space: pre-wrap;
                  font-size: 12px;
                  line-height: 1.6;
                  color: #e5e7eb;
                  overflow: auto;
                "
              >
{{ originalText || "(空)" }}
              </pre>
            </div>
            <div style="display: flex; flex-direction: column; min-width: 0">
              <div
                style="
                  font-size: 12px;
                  text-transform: uppercase;
                  letter-spacing: 0.06em;
                  color: #9ca3af;
                  margin-bottom: 6px;
                "
              >
                清洗后
              </div>
              <pre
                style="
                  flex: 1;
                  background: radial-gradient(
                    circle at top left,
                    rgba(15, 23, 42, 0.96),
                    rgba(15, 23, 42, 0.92)
                  );
                  padding: 12px;
                  border-radius: 12px;
                  border: 1px solid rgba(31, 41, 55, 0.95);
                  white-space: pre-wrap;
                  font-size: 12px;
                  line-height: 1.6;
                  color: #e5e7eb;
                  overflow: auto;
                "
              >
{{ sanitizedText || "(空)" }}
              </pre>
            </div>
          </div>
        </section>

        <section
          v-if="viewMode === 'report' && response"
          class="card card--subtle"
          style="flex: 1; overflow: hidden"
        >
          <ReportView :response="response" />
        </section>
      </section>

      <!-- 底部命中列表 -->
      <footer class="app-shell__footer">
        <FindingsList
          v-if="findings.length > 0"
          :findings="findings"
          :originalText="originalText"
          :onJump="(start, end) => {
            // 滚动到指定位置（简化实现）
            console.log('跳转到:', start, end);
          }"
        />
      </footer>
    </main>
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
