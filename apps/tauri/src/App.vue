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
              <div class="app-shell__title">PromptSanitizer</div>
              <div class="app-shell__subtitle">Privacy Security Lab</div>
            </div>
          </div>
        </div>

        <!-- Control Sections -->
        <div class="app-shell__controls">
          <!-- Primary Actions -->
          <div class="control-section">
            <div class="control-section__label">Operations</div>
            <button class="btn-action btn-action--primary" @click="handleSanitize">
              ⚡ Execute Sanitize
            </button>
            <button class="btn-action" @click="handleLoadFile" style="margin-top: 12px">
              📁 Load File
            </button>
          </div>

          <!-- Configuration Toggle -->
          <div class="control-section">
            <div class="control-section__label">Configuration</div>
            <button class="btn-action" @click="showConfig = !showConfig">
              {{ showConfig ? '⊗ Hide Config' : '⊕ Show Config' }}
            </button>
          </div>

          <!-- Config Panel (Collapsible) -->
          <div v-if="showConfig" class="control-section">
            <ConfigPanel :config="config" :onChange="setConfig" />
          </div>

          <!-- Statistics (if available) -->
          <div v-if="findings.length > 0" class="control-section">
            <div class="control-section__label">Statistics</div>
            <div style="padding: 12px; background: var(--color-bg-tertiary); border: 1px solid var(--color-border); font-family: var(--font-mono); font-size: 12px; line-height: 1.8;">
              <div style="display: flex; justify-content: space-between;">
                <span style="color: var(--color-text-muted);">Total Findings:</span>
                <span style="color: var(--color-accent); font-weight: 700;">{{ findings.length }}</span>
              </div>
              <div v-if="response" style="display: flex; justify-content: space-between; margin-top: 6px;">
                <span style="color: var(--color-text-muted);">Risk Score:</span>
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
            [Split]
          </button>
          <button
            class="view-switcher__btn"
            :class="{ 'view-switcher__btn--active': viewMode === 'diff' }"
            @click="viewMode = 'diff'"
          >
            [Diff]
          </button>
          <button
            class="view-switcher__btn"
            :class="{ 'view-switcher__btn--active': viewMode === 'report' }"
            @click="viewMode = 'report'"
          >
            [Report]
          </button>
        </div>

        <button class="theme-toggle" @click="toggleTheme">
          {{ theme === 'dark' ? '◐ DARK' : '◑ LIGHT' }}
        </button>
      </header>

      <!-- ============================================
           MAIN CONTENT AREA
           ============================================ -->
      <main class="app-shell__main">
        <!-- Split View: Original + Sanitized -->
        <template v-if="viewMode === 'split'">
          <div class="panel">
            <div class="panel__header">
              <h2 class="panel__title">◆ Original Text</h2>
            </div>
            <div class="panel__body">
              <MainPanel
                title="Original"
                :text="originalText"
                :onChange="setOriginalText"
                :findings="findings"
                :highlightMode="config.mode === 'annotate'"
              />
            </div>
          </div>

          <div class="panel">
            <div class="panel__header">
              <h2 class="panel__title">◆ Sanitized Output</h2>
            </div>
            <div class="panel__body">
              <MainPanel
                title="Sanitized"
                :text="sanitizedText"
                :readOnly="true"
                :findings="findings"
              />
            </div>
          </div>
        </template>

        <!-- Diff View -->
        <div v-if="viewMode === 'diff'" class="panel" style="flex: 1;">
          <div class="panel__header">
            <h2 class="panel__title">◆ Side-by-Side Comparison</h2>
          </div>
          <div class="panel__body">
            <div style="display: grid; grid-template-columns: 1fr 4px 1fr; gap: 20px; height: 100%;">
              <div style="display: flex; flex-direction: column;">
                <div style="font-size: 10px; text-transform: uppercase; letter-spacing: 0.15em; color: var(--color-text-muted); margin-bottom: 12px; font-weight: 700;">
                  ▸ Before
                </div>
                <pre class="text-editor" readonly style="flex: 1;">{{ originalText || '(Empty)' }}</pre>
              </div>
              
              <div style="width: 4px; background: linear-gradient(to bottom, transparent, var(--color-accent) 50%, transparent); opacity: 0.3;"></div>
              
              <div style="display: flex; flex-direction: column;">
                <div style="font-size: 10px; text-transform: uppercase; letter-spacing: 0.15em; color: var(--color-text-muted); margin-bottom: 12px; font-weight: 700;">
                  ▸ After
                </div>
                <pre class="text-editor" readonly style="flex: 1;">{{ sanitizedText || '(Empty)' }}</pre>
              </div>
            </div>
          </div>
        </div>

        <!-- Report View -->
        <div v-if="viewMode === 'report' && response" class="panel" style="flex: 1;">
          <div class="panel__header">
            <h2 class="panel__title">◆ Security Analysis Report</h2>
          </div>
          <div class="panel__body">
            <ReportView :response="response" />
          </div>
        </div>

        <!-- Empty State for Report -->
        <div v-if="viewMode === 'report' && !response" class="panel" style="flex: 1;">
          <div class="panel__body" style="display: flex; align-items: center; justify-content: center; flex-direction: column; gap: 16px; color: var(--color-text-muted);">
            <div style="font-size: 48px; opacity: 0.3;">⚠</div>
            <div style="font-family: var(--font-mono); font-size: 12px; text-transform: uppercase; letter-spacing: 0.1em;">
              No report available
            </div>
            <div style="font-size: 11px;">
              Execute sanitization to generate a security report
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
            <span>Detection Results</span>
            <span class="findings-count">{{ findings.length }}</span>
          </div>
          <select
            v-model="filterCategory"
            style="padding: 6px 12px; font-size: 11px; min-width: 140px;"
          >
            <option value="all">All Categories</option>
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
const theme = ref<"dark" | "light">("dark");
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

const applyTheme = (value: "dark" | "light") => {
  const body = document.body;
  body.classList.remove("theme-dark", "theme-light");
  body.classList.add(value === "dark" ? "theme-dark" : "theme-light");
};

const toggleTheme = () => {
  theme.value = theme.value === "dark" ? "light" : "dark";
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
  applyTheme(theme.value);
});

watch(
  theme,
  (value) => {
    applyTheme(value);
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
