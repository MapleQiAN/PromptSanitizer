<template>
  <div style="display: flex; flex-direction: column; gap: 24px; height: 100%; overflow-y: auto; padding: 4px;">
    <!-- Risk Score Hero -->
    <div
      style="
        padding: 32px;
        background: var(--color-bg-tertiary);
        border: var(--border-thick) solid var(--color-border-heavy);
        position: relative;
        overflow: hidden;
      "
      :style="{ borderColor: riskInfo.color }"
    >
      <!-- Diagonal decoration -->
      <div
        style="
          position: absolute;
          top: -50%;
          right: -10%;
          width: 300px;
          height: 300px;
          opacity: 0.05;
          transform: rotate(45deg);
        "
        :style="{ background: riskInfo.color }"
      />
      
      <div style="display: flex; align-items: center; gap: 32px; position: relative; z-index: 1;">
        <!-- Risk Score -->
        <div
          style="
            width: 120px;
            height: 120px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-family: var(--font-display);
            font-size: 56px;
            font-weight: 800;
            border: var(--border-thick) solid;
            clip-path: polygon(0 0, calc(100% - 16px) 0, 100% 16px, 100% 100%, 0 100%);
          "
          :style="{
            color: riskInfo.color,
            borderColor: riskInfo.color,
            boxShadow: `0 0 40px ${riskInfo.color}40`,
          }"
        >
          {{ response.risk_score }}
        </div>

        <!-- Risk Info -->
        <div style="flex: 1;">
          <h2
            style="
              font-size: 32px;
              font-weight: 800;
              text-transform: uppercase;
              letter-spacing: 0.05em;
              margin-bottom: 8px;
            "
            :style="{ color: riskInfo.color }"
          >
            {{ riskInfo.label }}
          </h2>
          <p style="color: var(--color-text-secondary); font-size: 13px; line-height: 1.6;">
            Comprehensive risk assessment based on detected patterns, sensitivity levels, and potential exposure vectors.
          </p>
        </div>
      </div>
    </div>

    <!-- Statistics Grid -->
    <div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px;">
      <div
        style="
          padding: 20px;
          background: var(--color-bg-tertiary);
          border: var(--border-thin) solid var(--color-border);
          clip-path: polygon(0 0, calc(100% - 10px) 0, 100% 10px, 100% 100%, 0 100%);
        "
      >
        <div style="font-size: 36px; font-weight: 800; color: var(--color-text-primary); font-family: var(--font-mono);">
          {{ response.stats.total_findings }}
        </div>
        <div style="font-size: 10px; text-transform: uppercase; letter-spacing: 0.15em; color: var(--color-text-muted); margin-top: 8px; font-weight: 700;">
          Total Findings
        </div>
      </div>

      <div
        style="
          padding: 20px;
          background: var(--color-bg-tertiary);
          border: var(--border-thin) solid var(--color-border);
          clip-path: polygon(0 0, calc(100% - 10px) 0, 100% 10px, 100% 100%, 0 100%);
        "
      >
        <div style="font-size: 36px; font-weight: 800; color: var(--color-risk-high); font-family: var(--font-mono);">
          {{ response.stats.high_risk_count }}
        </div>
        <div style="font-size: 10px; text-transform: uppercase; letter-spacing: 0.15em; color: var(--color-text-muted); margin-top: 8px; font-weight: 700;">
          High Risk
        </div>
      </div>

      <div
        style="
          padding: 20px;
          background: var(--color-bg-tertiary);
          border: var(--border-thin) solid var(--color-border);
          clip-path: polygon(0 0, calc(100% - 10px) 0, 100% 10px, 100% 100%, 0 100%);
        "
      >
        <div style="font-size: 36px; font-weight: 800; color: var(--color-risk-medium); font-family: var(--font-mono);">
          {{ response.stats.medium_risk_count }}
        </div>
        <div style="font-size: 10px; text-transform: uppercase; letter-spacing: 0.15em; color: var(--color-text-muted); margin-top: 8px; font-weight: 700;">
          Medium Risk
        </div>
      </div>

      <div
        style="
          padding: 20px;
          background: var(--color-bg-tertiary);
          border: var(--border-thin) solid var(--color-border);
          clip-path: polygon(0 0, calc(100% - 10px) 0, 100% 10px, 100% 100%, 0 100%);
        "
      >
        <div style="font-size: 36px; font-weight: 800; color: var(--color-risk-low); font-family: var(--font-mono);">
          {{ response.stats.low_risk_count }}
        </div>
        <div style="font-size: 10px; text-transform: uppercase; letter-spacing: 0.15em; color: var(--color-text-muted); margin-top: 8px; font-weight: 700;">
          Low Risk
        </div>
      </div>
    </div>

    <!-- Category Breakdown -->
    <div
      style="
        padding: 24px;
        background: var(--color-bg-tertiary);
        border: var(--border-normal) solid var(--color-border);
      "
    >
      <h3
        style="
          font-size: 12px;
          text-transform: uppercase;
          letter-spacing: 0.15em;
          color: var(--color-text-secondary);
          margin-bottom: 20px;
          font-weight: 800;
        "
      >
        ◆ Category Distribution
      </h3>
      <div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 12px;">
        <div
          v-for="[cat, count] in categoryEntries"
          :key="cat"
          style="
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 12px 16px;
            background: var(--color-bg-primary);
            border: var(--border-thin) solid var(--color-border);
            font-family: var(--font-mono);
            font-size: 12px;
            clip-path: polygon(0 0, calc(100% - 8px) 0, 100% 8px, 100% 100%, 0 100%);
          "
        >
          <span style="color: var(--color-text-secondary); text-transform: uppercase; letter-spacing: 0.05em;">
            {{ categoryLabels[cat] || cat }}
          </span>
          <span style="color: var(--color-accent); font-weight: 700; font-size: 14px;">
            {{ count }}
          </span>
        </div>
      </div>
    </div>

    <!-- Technical Details -->
    <div
      style="
        padding: 24px;
        background: var(--color-bg-tertiary);
        border: var(--border-normal) solid var(--color-border);
      "
    >
      <h3
        style="
          font-size: 12px;
          text-transform: uppercase;
          letter-spacing: 0.15em;
          color: var(--color-text-secondary);
          margin-bottom: 16px;
          font-weight: 800;
        "
      >
        ◆ Technical Report
      </h3>
      <div style="margin-bottom: 16px; font-family: var(--font-mono); font-size: 11px; color: var(--color-text-muted);">
        Engine Version: <span style="color: var(--color-accent); font-weight: 600;">{{ response.version }}</span>
      </div>
      <pre
        style="
          background: var(--color-bg-primary);
          border: var(--border-thin) solid var(--color-border);
          padding: 16px;
          overflow: auto;
          font-family: var(--font-mono);
          font-size: 11px;
          line-height: 1.5;
          max-height: 400px;
          color: var(--color-text-secondary);
        "
      >{{ JSON.stringify(response, null, 2) }}</pre>
    </div>

    <!-- Export Button -->
    <div style="display: flex; justify-content: flex-end;">
      <button class="btn-action btn-action--primary" @click="handleExport" style="padding: 12px 32px;">
        📥 Export Report (JSON)
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import type { Response } from "../types";

interface Props {
  response: Response;
}

const props = defineProps<Props>();

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

const getRiskLevel = (score: number) => {
  if (score >= 70) return { label: "High Risk", color: "var(--color-risk-high)" };
  if (score >= 40) return { label: "Medium Risk", color: "var(--color-risk-medium)" };
  return { label: "Low Risk", color: "var(--color-risk-low)" };
};

const riskInfo = computed(() => getRiskLevel(props.response.risk_score));

const categoryEntries = computed(() => {
  return Object.entries(props.response.stats.by_category);
});

const handleExport = () => {
  const safeReport = {
    ...props.response,
    findings: props.response.findings.map((f) => ({
      ...f,
    })),
  };

  const blob = new Blob([JSON.stringify(safeReport, null, 2)], {
    type: "application/json",
  });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `sanitization-report-${Date.now()}.json`;
  a.click();
  URL.revokeObjectURL(url);
};
</script>
