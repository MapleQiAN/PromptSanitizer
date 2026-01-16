<template>
  <div class="app-shell" :class="{ 'app-shell--sidebar-collapsed': sidebarCollapsed }">
    <!-- ============================================
         LEFT SIDEBAR - Control Panel
         ============================================ -->
    <aside class="app-shell__sidebar" :class="{ 'app-shell__sidebar--collapsed': sidebarCollapsed }">
      <!-- Toggle Button -->
      <button class="sidebar-toggle" @click="toggleSidebar" :title="sidebarCollapsed ? t.expandSidebar : t.collapseSidebar">
        <svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path
            v-if="!sidebarCollapsed"
            d="M10 12L6 8L10 4"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
          <path
            v-else
            d="M6 4L10 8L6 12"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
      </button>
      <!-- Brand / Logo -->
      <div class="app-shell__brand">
        <div class="app-shell__logo" v-if="!sidebarCollapsed">
          <div class="app-shell__logo-mark">
            <img src="/icon.png" alt="PromptSanitizer" class="app-shell__logo-img" />
          </div>
          <div class="app-shell__logo-text">
            <div class="app-shell__title">{{ t.brandTitle }}</div>
            <div class="app-shell__subtitle">{{ t.brandSubtitle }}</div>
          </div>
        </div>
        <!-- Collapsed Logo Icon -->
        <div v-else class="app-shell__logo-icon">
          <img src="/icon.png" alt="PromptSanitizer" class="app-shell__logo-img" />
        </div>
      </div>

      <!-- Control Sections -->
      <div class="app-shell__controls" v-show="!sidebarCollapsed">
        <!-- Config Panel (Always Visible) -->
        <div class="control-section">
          <div class="control-section__label">{{ t.configuration }}</div>
          <ConfigPanel :config="config" :onChange="setConfig" :lang="lang" />
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
      <!-- Split View: Original + Sanitized (Text or Image) -->
      <template v-if="viewMode === 'split'">
        <div class="panel">
          <div class="panel__header">
            <h2 class="panel__title">{{ t.originalText }}</h2>
            <div class="panel__header-actions" style="display: flex; gap: 12px; align-items: center; flex-wrap: wrap;">
              <div style="display: flex; gap: 8px; align-items: center; padding-right: 16px; margin-right: 8px; border-right: 2px solid var(--color-border-light);">
                <button class="btn-action" @click="handleLoadFile" style="width: auto; white-space: nowrap; padding: 8px 16px;">
                {{ t.loadFile }}
              </button>
                <button class="btn-action" @click="handleLoadImage" style="width: auto; white-space: nowrap; padding: 8px 16px;">
                  {{ lang === 'zh' ? '📷 加载图片' : '📷 Load Image' }}
                </button>
              </div>
              <button 
                v-if="!isImageMode"
                class="btn-action btn-action--primary" 
                @click="handleSanitize" 
                style="width: auto; white-space: nowrap; padding: 8px 20px;"
              >
                {{ t.executeSanitize }}
              </button>
            </div>
          </div>
          <div class="panel__body">
            <!-- Text Mode -->
            <MainPanel
              v-if="!isImageMode"
              title="Original"
              :text="originalText"
              :onChange="setOriginalText"
              :findings="findings"
              :highlightMode="config.mode === 'annotate'"
              :fontSize="config.fontSize || 16"
            />
            <!-- Image Mode -->
            <ImagePanel
              v-else
              ref="imagePanelRef"
              :lang="lang"
              @maskedImage="handleMaskedImage"
            />
          </div>
        </div>

        <div class="panel">
          <div class="panel__header">
            <h2 class="panel__title">{{ t.sanitizedOutput }}</h2>
            <div class="panel__header-actions"></div>
          </div>
          <div class="panel__body">
            <!-- Text Mode -->
            <MainPanel
              v-if="!isImageMode"
              title="Sanitized"
              :text="sanitizedText"
              :readOnly="true"
              :findings="findings"
              :fontSize="config.fontSize || 16"
            />
            <!-- Image Mode -->
            <div 
              v-else 
              style="display: flex; flex-direction: column; height: 100%; align-items: center; justify-content: center; padding: 24px; overflow: auto; transition: all 0.3s ease;"
              :style="isRightPanelDragging ? { backgroundColor: 'var(--color-bg-secondary)' } : {}"
              @dragover.prevent="handleRightPanelDragOver"
              @dragleave.prevent="handleRightPanelDragLeave"
              @drop.prevent="handleRightPanelDrop"
            >
              <div v-if="!maskedImageUrl" class="empty-state" style="text-align: center; color: var(--color-text-muted); padding: 40px 20px; width: 100%;">
                <div style="font-size: 64px; opacity: 0.15; margin-bottom: 16px; transition: transform 0.3s ease;" :style="isRightPanelDragging ? { transform: 'scale(1.1)', opacity: '0.3' } : {}">✨</div>
                <div style="font-family: var(--font-display); font-size: 15px; font-weight: 600; margin-bottom: 8px;">
                  {{ isRightPanelDragging ? (lang === 'zh' ? '松开鼠标以加载图片' : 'Release to load image') : (lang === 'zh' ? '等待打码处理或拖拽图片' : 'Waiting for masking or drag image') }}
                </div>
                <div style="font-size: 13px; opacity: 0.7;">
                  {{ lang === 'zh' ? '在左侧检测文本并应用打码后，结果将显示在这里' : 'After detecting text and applying mask on the left, the result will appear here' }}
                </div>
              </div>
              <img
                v-else
                :src="maskedImageUrl"
                alt="Masked"
                style="max-width: 100%; max-height: 100%; object-fit: contain; border-radius: var(--radius-sm); box-shadow: 0 4px 12px rgba(0,0,0,0.1);"
              />
            </div>
          </div>
        </div>
      </template>

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
        <CustomSelect
          v-model="filterCategory"
          :options="[
            { value: 'all', label: t.allCategories },
            ...categories.map((cat) => ({
              value: cat,
              label: categoryLabels[cat] || cat,
            })),
          ]"
          style="min-width: 160px;"
        />
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
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from "vue";
import { useMessage } from "naive-ui";
import MainPanel from "./MainPanel.vue";
import ConfigPanel from "./ConfigPanel.vue";
import FindingsList from "./FindingsList.vue";
import ReportView from "./ReportView.vue";
import ImagePanel from "./ImagePanel.vue";
import CustomSelect from "./CustomSelect.vue";
import type { Config, Finding, Response } from "../types";

const message = useMessage();

// 动态导入 Tauri API，避免在非 Tauri 环境中出错
let invokeFn: typeof import("@tauri-apps/api/core").invoke | null = null;
let initPromise: Promise<void> | null = null;

// 初始化 Tauri API（只初始化一次，不阻塞渲染）
const initTauri = async () => {
  if (initPromise) {
    return initPromise;
  }
  
  initPromise = (async () => {
    try {
      if (typeof window !== "undefined") {
        // 使用动态导入，不阻塞主线程
        const tauriApi = await import("@tauri-apps/api/core");
        invokeFn = tauriApi.invoke;
      }
    } catch (error) {
      console.warn("Tauri API not available:", error);
    }
  })();
  
  // 不等待完成，让 UI 先渲染
  return initPromise;
};

// 安全的 invoke 包装函数
const invoke = async <T = any>(cmd: string, args?: any): Promise<T> => {
  // 如果还没初始化，先初始化（但不会阻塞）
  if (!invokeFn && !initPromise) {
    initTauri();
  }
  // 等待初始化完成
  if (!invokeFn) {
    await initTauri();
  }
  if (!invokeFn) {
    throw new Error("Tauri API is not available. Please run the app using 'npm run tauri:dev' instead of 'npm run dev'");
  }
  return await invokeFn<T>(cmd, args);
};

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
    collapseSidebar: "收起侧边栏",
    expandSidebar: "展开侧边栏",
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
    collapseSidebar: "Collapse Sidebar",
    expandSidebar: "Expand Sidebar",
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
// 配置面板始终显示，不再需要 showConfig
const viewMode = ref<"split" | "report">("split");
const isImageMode = ref(false);
const maskedImageUrl = ref<string>("");
const isRightPanelDragging = ref(false);
const imagePanelRef = ref<InstanceType<typeof ImagePanel> | null>(null);
const theme = ref<"dark" | "light">("light");
const filterCategory = ref<string>("all");
const sidebarCollapsed = ref(false);

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value;
};

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

const getRiskColor = (score: number) => {
  if (score >= 70) return "var(--color-risk-high)";
  if (score >= 40) return "var(--color-risk-medium)";
  return "var(--color-risk-low)";
};

onMounted(() => {
  // 初始化主题（同步操作，立即执行）
  const prefersDark = window.matchMedia?.("(prefers-color-scheme: dark)")?.matches;
  theme.value = prefersDark ? "dark" : "light";
  applyTheme(theme.value, themeStyle.value);
  
  // 异步初始化 Tauri API（不阻塞渲染）
  // 使用 requestIdleCallback 或 setTimeout 延迟执行，确保 UI 先渲染
  if (window.requestIdleCallback) {
    window.requestIdleCallback(() => {
      initTauri();
    });
  } else {
    setTimeout(() => {
      initTauri();
    }, 0);
  }
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
  if (isImageMode.value) {
    message.info(lang.value === "zh" ? "图片模式下请使用左侧的检测和打码功能" : "In image mode, please use the detect and mask functions on the left");
    return;
  }

  if (!originalText.value.trim()) {
    message.warning(lang.value === "zh" ? "请输入要净化的文本" : "Please input text to sanitize");
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
    message.error(lang.value === "zh" ? `净化失败: ${error}` : `Sanitization failed: ${error}`);
  }
};

const handleLoadFile = async () => {
  try {
    isImageMode.value = false;
    const filePath = await invoke<string>("open_file_dialog");
    if (filePath) {
      const content = await invoke<string>("read_file", { path: filePath });
      originalText.value = content;
      message.success(lang.value === "zh" ? "文件加载成功" : "File loaded successfully");
    }
  } catch (error) {
    console.error("Failed to load file:", error);
    message.error(lang.value === "zh" ? `文件加载失败: ${error}` : `Failed to load file: ${error}`);
  }
};

const handleLoadImage = () => {
  isImageMode.value = true;
  originalText.value = "";
  sanitizedText.value = "";
  findings.value = [];
  response.value = null;
  maskedImageUrl.value = "";
};

const handleMaskedImage = (url: string) => {
  maskedImageUrl.value = url;
};

const handleRightPanelDragOver = (event: DragEvent) => {
  event.preventDefault();
  if (!isRightPanelDragging.value) {
    isRightPanelDragging.value = true;
  }
};

const handleRightPanelDragLeave = (event: DragEvent) => {
  event.preventDefault();
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect();
  const x = event.clientX;
  const y = event.clientY;
  
  if (x < rect.left || x > rect.right || y < rect.top || y > rect.bottom) {
    isRightPanelDragging.value = false;
  }
};

const handleRightPanelDrop = (event: DragEvent) => {
  event.preventDefault();
  isRightPanelDragging.value = false;
  
  const files = event.dataTransfer?.files;
  if (files && files.length > 0) {
    const file = files[0];
    if (file.type.startsWith("image/")) {
      // 切换到图片模式
      if (!isImageMode.value) {
        handleLoadImage();
      }
      // 使用 ImagePanel 的方法来处理文件
      setTimeout(() => {
        if (imagePanelRef.value && (imagePanelRef.value as any).processImageFile) {
          (imagePanelRef.value as any).processImageFile(file);
        }
      }, 100);
    } else {
      message.warning(lang.value === "zh" ? "请拖入有效的图片文件" : "Please drag a valid image file");
    }
  }
};
</script>
