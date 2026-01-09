<template>
  <n-config-provider :theme="naiveTheme">
    <div class="app-shell">
      <!-- ============================================
           LEFT SIDEBAR - Control Panel
           ============================================ -->
      <aside class="app-shell__sidebar">
        <!-- Brand / Logo -->
        <div class="app-shell__brand">
          <div class="app-shell__logo">
            <div class="app-shell__logo-mark">PS</div>
            <div class="app-shell__logo-text">
              <div class="app-shell__title">{{ t.brandTitle }}</div>
              <div class="app-shell__subtitle">{{ t.brandSubtitle }}</div>
            </div>
          </div>
        </div>

        <!-- Control Sections -->
        <div class="app-shell__controls">
          <!-- Primary Actions -->
          <div class="control-section">
            <div class="control-section__label">{{ t.operations }}</div>
            <button class="btn-action btn-action--primary" @click="handleSanitize">
              {{ t.executeSanitize }}
            </button>
            <button class="btn-action" @click="handleLoadFile" style="margin-top: 12px">
              {{ t.loadFile }}
            </button>
          </div>

          <!-- Configuration Toggle -->
          <div class="control-section">
            <div class="control-section__label">{{ t.configuration }}</div>
            <button class="btn-action" @click="showConfig = !showConfig">
              {{ showConfig ? t.hideConfig : t.showConfig }}
            </button>
          </div>

          <!-- Config Panel (Collapsible) -->
          <div v-if="showConfig" class="control-section">
            <ConfigPanel :config="config" :onChange="setConfig" />
          </div>

          <!-- Statistics (if available) -->
          <div v-if="findings.length > 0" class="control-section">
            <div class="control-section__label">{{ t.statistics }}</div>
            <div style="padding: 16px; background: var(--color-bg-tertiary); border: 3px solid var(--color-border); border-radius: var(--radius-md); font-family: var(--font-mono); font-size: 13px; line-height: 2;">
              <div style="display: flex; justify-content: space-between;">
                <span style="color: var(--color-text-secondary);">{{ t.totalFindings }}</span>
                <span style="color: var(--color-primary); font-weight: 700;">{{ findings.length }}</span>
              </div>
              <div v-if="response" style="display: flex; justify-content: space-between; margin-top: 8px;">
                <span style="color: var(--color-text-secondary);">{{ t.riskScore }}</span>
                <span :style="{ color: getRiskColor(response.risk_score), fontWeight: 700 }">
                  {{ response.risk_score }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </aside>

      <!-- ============================================
           HEADER BAR - View Controls
           ============================================ -->
      <header class="app-shell__header">
        <div class="view-switcher">
          <button
            class="view-switcher__btn"
            :class="{ 'view-switcher__btn--active': viewMode === 'split' }"
            @click="viewMode = 'split'"
          >
            {{ lang === 'zh' ? '分栏' : 'Split' }}
          </button>
          <button
            class="view-switcher__btn"
            :class="{ 'view-switcher__btn--active': viewMode === 'diff' }"
            @click="viewMode = 'diff'"
          >
            {{ lang === 'zh' ? '对比' : 'Diff' }}
          </button>
          <button
            class="view-switcher__btn"
            :class="{ 'view-switcher__btn--active': viewMode === 'report' }"
            @click="viewMode = 'report'"
          >
            {{ lang === 'zh' ? '报告' : 'Report' }}
          </button>
        </div>

        <div style="display: flex; gap: 12px; align-items: center;">
          <button class="style-toggle" @click="toggleStyle">
            {{ themeStyle === 'kawaii' ? t.kawaiiStyle : t.horizonStyle }}
          </button>
          <button class="lang-toggle" @click="toggleLang">
            {{ lang === 'zh' ? '🌏 中文' : '🌍 EN' }}
          </button>
          <button class="theme-toggle" @click="toggleTheme">
            {{ theme === 'dark' ? t.darkMode : t.lightMode }}
          </button>
        </div>
      </header>

      <!-- ============================================
           MAIN CONTENT AREA
           ============================================ -->
      <main class="app-shell__main">
        <!-- Split View: Original + Sanitized -->
        <template v-if="viewMode === 'split'">
          <div class="panel">
            <div class="panel__header">
              <h2 class="panel__title">{{ t.originalText }}</h2>
            </div>
            <div class="panel__body">
              <MainPanel
                title="Original"
                :text="originalText"
                :onChange="setOriginalText"
                :findings="findings"
                :highlightMode="config.mode === 'annotate'"
                :fontSize="config.fontSize || 16"
              />
            </div>
          </div>

          <div class="panel">
            <div class="panel__header">
              <h2 class="panel__title">{{ t.sanitizedOutput }}</h2>
            </div>
            <div class="panel__body">
              <MainPanel
                title="Sanitized"
                :text="sanitizedText"
                :readOnly="true"
                :findings="findings"
                :fontSize="config.fontSize || 16"
              />
            </div>
          </div>
        </template>

        <!-- Diff View -->
        <div v-if="viewMode === 'diff'" class="panel" style="flex: 1;">
          <div class="panel__header">
            <h2 class="panel__title">{{ t.comparison }}</h2>
          </div>
          <div class="panel__body">
            <div style="display: grid; grid-template-columns: 1fr 4px 1fr; gap: 20px; height: 100%;">
              <div style="display: flex; flex-direction: column;">
                <div style="font-size: 13px; font-weight: 600; color: var(--color-text-secondary); margin-bottom: 12px; font-family: var(--font-display);">
                  {{ t.before }}
                </div>
                <pre class="text-editor" readonly :style="{ flex: 1, fontSize: `${config.fontSize || 16}px` }">{{ originalText || `(${lang === 'zh' ? '空' : 'Empty'})` }}</pre>
              </div>
              
              <div style="width: 4px; background: linear-gradient(to bottom, transparent, var(--color-primary) 50%, transparent); opacity: 0.3; border-radius: var(--radius-full);"></div>
              
              <div style="display: flex; flex-direction: column;">
                <div style="font-size: 13px; font-weight: 600; color: var(--color-text-secondary); margin-bottom: 12px; font-family: var(--font-display);">
                  {{ t.after }}
                </div>
                <pre class="text-editor" readonly :style="{ flex: 1, fontSize: `${config.fontSize || 16}px` }">{{ sanitizedText || `(${lang === 'zh' ? '空' : 'Empty'})` }}</pre>
              </div>
            </div>
          </div>
        </div>

        <!-- Report View -->
        <div v-if="viewMode === 'report' && response" class="panel" style="flex: 1;">
          <div class="panel__header">
            <h2 class="panel__title">{{ t.securityReport }}</h2>
          </div>
          <div class="panel__body">
            <ReportView :response="response" />
          </div>
        </div>

        <!-- Empty State for Report -->
        <div v-if="viewMode === 'report' && !response" class="panel" style="flex: 1;">
          <div class="panel__body" style="display: flex; align-items: center; justify-content: center; flex-direction: column; gap: 20px; color: var(--color-text-muted);">
            <div style="font-size: 64px; opacity: 0.2;">🎀</div>
            <div style="font-family: var(--font-display); font-size: 18px; font-weight: 600;">
              {{ t.noReport }}
            </div>
            <div style="font-size: 14px;">
              {{ t.noReportDesc }}
            </div>
          </div>
        </div>
      </main>

      <!-- ============================================
           FOOTER - Findings List
           ============================================ -->
      <footer v-if="findings.length > 0" class="app-shell__footer">
        <div class="findings-header">
          <div class="findings-title">
            <span>{{ t.detectionResults }}</span>
            <span class="findings-count">{{ findings.length }}</span>
          </div>
          <select
            v-model="filterCategory"
            style="padding: 10px 16px; font-size: 13px; min-width: 160px;"
          >
            <option value="all">{{ t.allCategories }}</option>
            <option v-for="cat in categories" :key="cat" :value="cat">
              {{ categoryLabels[cat] || cat }}
            </option>
          </select>
        </div>
        <div class="findings-list">
          <FindingsList
            :findings="filteredFindings"
            :originalText="originalText"
            :onJump="(start, end) => console.log('Jump to:', start, end)"
          />
        </div>
      </footer>
    </div>
  </n-config-provider>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from "vue";
import { darkTheme, type GlobalTheme } from "naive-ui";
import { invoke } from "@tauri-apps/api/core";
import MainPanel from "./components/MainPanel.vue";
import ConfigPanel from "./components/ConfigPanel.vue";
import FindingsList from "./components/FindingsList.vue";
import ReportView from "./components/ReportView.vue";
import type { Config, Finding, Response } from "./types";

// 语言配置
type Language = "zh" | "en";
const lang = ref<Language>("zh");

// 主题风格配置
type ThemeStyle = "kawaii" | "horizon";
const themeStyle = ref<ThemeStyle>("kawaii");

const i18n = {
  zh: {
    brandTitle: "提示词净化器",
    brandSubtitle: "隐私安全实验室",
    operations: "操作",
    executeSanitize: "⚡ 执行净化",
    loadFile: "📁 加载文件",
    configuration: "配置",
    hideConfig: "⊗ 隐藏配置",
    showConfig: "⊕ 显示配置",
    statistics: "统计信息",
    totalFindings: "发现总数:",
    riskScore: "风险评分:",
    originalText: "💝 原始文本",
    sanitizedOutput: "✨ 净化输出",
    comparison: "🎀 对比视图",
    securityReport: "🦋 安全报告",
    before: "▸ 净化前",
    after: "▸ 净化后",
    noReport: "暂无报告",
    noReportDesc: "执行净化以生成安全报告",
    detectionResults: "检测结果",
    allCategories: "所有类别",
    inputPlaceholder: "请输入要净化的文本...",
    darkMode: "🌙 深色",
    lightMode: "☀️ 浅色",
    kawaiiStyle: "🎀 可爱风",
    horizonStyle: "🚀 科技风",
  },
  en: {
    brandTitle: "PromptSanitizer",
    brandSubtitle: "Privacy Security Lab",
    operations: "Operations",
    executeSanitize: "⚡ Execute Sanitize",
    loadFile: "📁 Load File",
    configuration: "Configuration",
    hideConfig: "⊗ Hide Config",
    showConfig: "⊕ Show Config",
    statistics: "Statistics",
    totalFindings: "Total Findings:",
    riskScore: "Risk Score:",
    originalText: "💝 Original Text",
    sanitizedOutput: "✨ Sanitized Output",
    comparison: "🎀 Comparison",
    securityReport: "🦋 Security Report",
    before: "▸ Before",
    after: "▸ After",
    noReport: "No Report Available",
    noReportDesc: "Execute sanitization to generate a security report",
    detectionResults: "Detection Results",
    allCategories: "All Categories",
    inputPlaceholder: "Enter text to sanitize...",
    darkMode: "🌙 Dark",
    lightMode: "☀️ Light",
    kawaiiStyle: "🎀 Kawaii",
    horizonStyle: "🚀 Horizon",
  },
};

const t = computed(() => i18n[lang.value]);

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
  fontSize: 16,
});
const showConfig = ref(false);
const viewMode = ref<"split" | "diff" | "report">("split");
const theme = ref<"dark" | "light">("light");
const filterCategory = ref<string>("all");

const categoryLabels: Record<string, string> = {
  phone: "Phone",
  email: "Email",
  id_card: "ID Card",
  ip: "IP Address",
  domain: "Domain/URL",
  token: "Token/Key",
  password: "Password",
  private_key: "Private Key",
};

const categories = computed(() => {
  return Array.from(new Set(findings.value.map((f) => f.type)));
});

const filteredFindings = computed(() => {
  return filterCategory.value === "all"
    ? findings.value
    : findings.value.filter((f) => f.type === filterCategory.value);
});

const applyTheme = (colorMode: "dark" | "light", style: ThemeStyle) => {
  const body = document.body;
  body.classList.remove("theme-dark", "theme-light", "style-kawaii", "style-horizon");
  body.classList.add(colorMode === "dark" ? "theme-dark" : "theme-light");
  body.classList.add(style === "kawaii" ? "style-kawaii" : "style-horizon");
};

const toggleTheme = () => {
  theme.value = theme.value === "dark" ? "light" : "dark";
};

const toggleLang = () => {
  lang.value = lang.value === "zh" ? "en" : "zh";
};

const toggleStyle = () => {
  themeStyle.value = themeStyle.value === "kawaii" ? "horizon" : "kawaii";
};

const naiveTheme = computed<GlobalTheme | null>(() =>
  theme.value === "dark" ? darkTheme : null
);

const getRiskColor = (score: number) => {
  if (score >= 70) return "var(--color-risk-high)";
  if (score >= 40) return "var(--color-risk-medium)";
  return "var(--color-risk-low)";
};

onMounted(() => {
  const prefersDark = window.matchMedia?.("(prefers-color-scheme: dark)")?.matches;
  theme.value = prefersDark ? "dark" : "light";
  applyTheme(theme.value, themeStyle.value);
});

watch(
  [theme, themeStyle],
  ([newTheme, newStyle]) => {
    applyTheme(newTheme, newStyle);
  },
  { flush: "post" }
);

const setOriginalText = (text: string) => {
  originalText.value = text;
};

const setConfig = (newConfig: Config) => {
  config.value = newConfig;
};

const handleSanitize = async () => {
  if (!originalText.value.trim()) {
    alert("Please input text to sanitize");
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
    console.error("Sanitization failed:", error);
    alert(`Sanitization failed: ${error}`);
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
    console.error("Failed to load file:", error);
  }
};
</script>
