<template>
  <div style="display: flex; flex-direction: column; gap: 24px; height: 100%; overflow-y: auto; padding: 4px;">
    <!-- Risk Score Hero -->
    <div
      style="
        padding: 40px;
        background: linear-gradient(135deg, var(--color-bg-tertiary), var(--color-bg-secondary));
        border: 4px solid;
        border-radius: var(--radius-lg);
        position: relative;
        overflow: hidden;
        box-shadow: var(--shadow-lg);
      "
      :style="{ borderColor: riskInfo.color }"
    >
      <!-- Cute bubble decoration -->
      <div
        style="
          position: absolute;
          top: -80px;
          right: -80px;
          width: 200px;
          height: 200px;
          opacity: 0.08;
          border-radius: var(--radius-full);
          animation: bubble-float-1 6s ease-in-out infinite;
        "
        :style="{ background: riskInfo.color }"
      />
      
      <div style="display: flex; align-items: center; gap: 40px; position: relative; z-index: 1;">
        <!-- Risk Score -->
        <div
          style="
            width: 140px;
            height: 140px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-family: var(--font-display);
            font-size: 64px;
            font-weight: 700;
            border: 5px solid;
            border-radius: var(--radius-lg);
            animation: pulse-glow 3s ease-in-out infinite;
          "
          :style="{
            color: riskInfo.color,
            borderColor: riskInfo.color,
            boxShadow: `0 8px 32px ${riskInfo.color}40`,
          }"
        >
          {{ response.risk_score }}
        </div>

        <!-- Risk Info -->
        <div style="flex: 1;">
          <h2
            style="
              font-size: 36px;
              font-weight: 700;
              margin-bottom: 12px;
              font-family: var(--font-display);
            "
            :style="{ color: riskInfo.color }"
          >
            {{ riskInfo.label }}
          </h2>
          <p style="color: var(--color-text-secondary); font-size: 15px; line-height: 1.7;">
            Comprehensive risk assessment based on detected patterns, sensitivity levels, and potential exposure vectors.
          </p>
        </div>
      </div>
    </div>

    <!-- Statistics Grid -->
    <div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px;">
      <div
        style="
          padding: 28px;
          background: var(--color-bg-tertiary);
          border: 3px solid var(--color-border);
          border-radius: var(--radius-md);
          box-shadow: var(--shadow-sm);
          transition: all 0.3s var(--ease-smooth);
        "
        class="stat-card"
      >
        <div style="font-size: 48px; font-weight: 700; color: var(--color-text-primary); font-family: var(--font-display);">
          {{ response.stats.total_findings }}
        </div>
        <div style="font-size: 13px; color: var(--color-text-secondary); margin-top: 10px; font-weight: 600; font-family: var(--font-display);">
          Total Findings
        </div>
      </div>

      <div
        style="
          padding: 28px;
          background: var(--color-bg-tertiary);
          border: 3px solid var(--color-border);
          border-radius: var(--radius-md);
          box-shadow: var(--shadow-sm);
          transition: all 0.3s var(--ease-smooth);
        "
        class="stat-card"
      >
        <div style="font-size: 48px; font-weight: 700; color: var(--color-risk-high); font-family: var(--font-display);">
          {{ response.stats.high_risk_count }}
        </div>
        <div style="font-size: 13px; color: var(--color-text-secondary); margin-top: 10px; font-weight: 600; font-family: var(--font-display);">
          High Risk
        </div>
      </div>

      <div
        style="
          padding: 28px;
          background: var(--color-bg-tertiary);
          border: 3px solid var(--color-border);
          border-radius: var(--radius-md);
          box-shadow: var(--shadow-sm);
          transition: all 0.3s var(--ease-smooth);
        "
        class="stat-card"
      >
        <div style="font-size: 48px; font-weight: 700; color: var(--color-risk-medium); font-family: var(--font-display);">
          {{ response.stats.medium_risk_count }}
        </div>
        <div style="font-size: 13px; color: var(--color-text-secondary); margin-top: 10px; font-weight: 600; font-family: var(--font-display);">
          Medium Risk
        </div>
      </div>

      <div
        style="
          padding: 28px;
          background: var(--color-bg-tertiary);
          border: 3px solid var(--color-border);
          border-radius: var(--radius-md);
          box-shadow: var(--shadow-sm);
          transition: all 0.3s var(--ease-smooth);
        "
        class="stat-card"
      >
        <div style="font-size: 48px; font-weight: 700; color: var(--color-risk-low); font-family: var(--font-display);">
          {{ response.stats.low_risk_count }}
        </div>
        <div style="font-size: 13px; color: var(--color-text-secondary); margin-top: 10px; font-weight: 600; font-family: var(--font-display);">
          Low Risk
        </div>
      </div>
    </div>

    <!-- Category Breakdown -->
    <div
      style="
        padding: 32px;
        background: var(--color-bg-tertiary);
        border: 3px solid var(--color-border);
        border-radius: var(--radius-lg);
        box-shadow: var(--shadow-sm);
      "
    >
      <h3
        style="
          font-size: 18px;
          color: var(--color-text-primary);
          margin-bottom: 24px;
          font-weight: 700;
          font-family: var(--font-display);
          display: flex;
          align-items: center;
          gap: 10px;
        "
      >
        <span style="font-size: 20px;">📊</span>
        <span>Category Distribution</span>
      </h3>
      <div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 16px;">
        <div
          v-for="[cat, count] in categoryEntries"
          :key="cat"
          style="
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 16px 20px;
            background: var(--color-bg-secondary);
            border: 3px solid var(--color-border);
            border-radius: var(--radius-md);
            font-family: var(--font-display);
            font-size: 14px;
            transition: all 0.3s var(--ease-smooth);
          "
          class="category-item"
        >
          <span style="color: var(--color-text-secondary); font-weight: 600;">
            {{ categoryLabels[cat] || cat }}
          </span>
          <span style="color: var(--color-primary); font-weight: 700; font-size: 18px;">
            {{ count }}
          </span>
        </div>
      </div>
    </div>

    <!-- Technical Details -->
    <div
      style="
        padding: 32px;
        background: var(--color-bg-tertiary);
        border: 3px solid var(--color-border);
        border-radius: var(--radius-lg);
        box-shadow: var(--shadow-sm);
      "
    >
      <h3
        style="
          font-size: 18px;
          color: var(--color-text-primary);
          margin-bottom: 20px;
          font-weight: 700;
          font-family: var(--font-display);
          display: flex;
          align-items: center;
          gap: 10px;
        "
      >
        <span style="font-size: 20px;">🔧</span>
        <span>Technical Report</span>
      </h3>
      <div style="margin-bottom: 20px; font-family: var(--font-display); font-size: 14px; color: var(--color-text-secondary); font-weight: 600;">
        Engine Version: <span style="color: var(--color-primary); font-weight: 700;">{{ response.version }}</span>
      </div>
      <pre
        style="
          background: var(--color-bg-secondary);
          border: 3px solid var(--color-border);
          border-radius: var(--radius-md);
          padding: 20px;
          overflow: auto;
          font-family: var(--font-mono);
          font-size: 12px;
          line-height: 1.7;
          max-height: 400px;
          color: var(--color-text-secondary);
        "
      >{{ JSON.stringify(response, null, 2) }}</pre>
    </div>

    <!-- Export Button -->
    <div style="display: flex; justify-content: center;">
      <button class="btn-action btn-action--primary" @click="handleExport" style="padding: 16px 48px; font-size: 16px;">
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

<style scoped>
.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}

.category-item:hover {
  transform: translateX(4px);
  border-color: var(--color-primary);
  box-shadow: var(--shadow-sm);
}
</style>
