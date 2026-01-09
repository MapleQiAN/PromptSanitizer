<template>
  <div style="flex: 1; padding: 18px 18px 20px; overflow: auto">
    <div
      style="
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 18px;
      "
    >
      <div>
        <h2 style="font-size: 18px; margin-bottom: 4px">清洗报告</h2>
        <div style="font-size: 12px; color: #9ca3af">
          总览本次文本脱敏的整体风险与命中分布
        </div>
      </div>
      <button class="btn-ghost" @click="handleExport" style="padding-inline: 14px">
        导出 JSON
      </button>
    </div>

    <!-- 风险评分 -->
    <div
      class="card"
      style="
        margin-bottom: 16px;
        padding: 16px 18px;
        border-width: 2px;
      "
      :style="{ borderColor: riskInfo.color }"
    >
      <div style="display: flex; align-items: center; gap: 18px">
        <div
          :style="{
            width: '72px',
            height: '72px',
            borderRadius: '999px',
            border: `2px solid ${riskInfo.color}`,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            fontSize: '26px',
            fontWeight: '700',
            color: riskInfo.color,
            boxShadow: '0 0 30px rgba(15,23,42,0.9)',
          }"
        >
          {{ response.risk_score }}
        </div>
        <div>
          <div style="font-size: 18px; font-weight: 600; margin-bottom: 4px">
            {{ riskInfo.label }}
          </div>
          <div style="color: #9ca3af; font-size: 13px">
            根据当前策略综合评估得到的文本泄露风险评分
          </div>
        </div>
      </div>
    </div>

    <!-- 统计信息 -->
    <div
      class="card card--subtle"
      style="
        padding: 16px 18px 18px;
        margin-bottom: 16px;
      "
    >
      <h3 style="margin-bottom: 14px; font-size: 14px">统计信息</h3>
      <div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 14px">
        <div>
          <div style="font-size: 22px; font-weight: 600">
            {{ response.stats.total_findings }}
          </div>
          <div style="color: #9ca3af; font-size: 12px">总命中数</div>
        </div>
        <div>
          <div style="font-size: 22px; font-weight: 600; color: #f97373">
            {{ response.stats.high_risk_count }}
          </div>
          <div style="color: #9ca3af; font-size: 12px">高危</div>
        </div>
        <div>
          <div style="font-size: 22px; font-weight: 600; color: #fb923c">
            {{ response.stats.medium_risk_count }}
          </div>
          <div style="color: #9ca3af; font-size: 12px">中危</div>
        </div>
        <div>
          <div style="font-size: 22px; font-weight: 600; color: #4ade80">
            {{ response.stats.low_risk_count }}
          </div>
          <div style="color: #9ca3af; font-size: 12px">低危</div>
        </div>
      </div>

      <div style="margin-top: 16px">
        <h4 style="margin-bottom: 10px; font-size: 13px">按类别统计</h4>
        <div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 8px">
          <div
            v-for="[cat, count] in categoryEntries"
            :key="cat"
            style="
              display: flex;
              justify-content: space-between;
              padding: 8px 10px;
              border-radius: 999px;
              background: rgba(15, 23, 42, 0.92);
              border: 1px solid rgba(31, 41, 55, 0.95);
              font-size: 12px;
            "
          >
            <span>{{ categoryLabels[cat] || cat }}</span>
            <span style="font-weight: 600">{{ count }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 详细信息 -->
    <div
      class="card card--subtle"
      style="
        padding: 16px 18px 18px;
      "
    >
      <h3 style="margin-bottom: 12px; font-size: 14px">详细信息</h3>
      <div style="font-size: 11px; color: #9ca3af; margin-bottom: 12px">
        引擎版本: {{ response.version }}
      </div>
      <pre
        style="
          background: radial-gradient(
            circle at top left,
            rgba(15, 23, 42, 0.96),
            rgba(15, 23, 42, 0.92)
          );
          padding: 12px 14px;
          border-radius: 12px;
          border: 1px solid rgba(31, 41, 55, 0.95);
          overflow: auto;
          font-size: 11px;
          max-height: 380px;
        "
      >
{{ JSON.stringify(response, null, 2) }}
      </pre>
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
  phone: "手机号",
  email: "邮箱",
  id_card: "身份证",
  ip: "IP地址",
  domain: "域名/URL",
  token: "Token/Key",
  password: "密码",
  private_key: "私钥",
};

const getRiskLevel = (score: number) => {
  if (score >= 70) return { label: "高危", color: "#f44336" };
  if (score >= 40) return { label: "中危", color: "#ff9800" };
  return { label: "低危", color: "#4caf50" };
};

const riskInfo = computed(() => getRiskLevel(props.response.risk_score));

const categoryEntries = computed(() => {
  return Object.entries(props.response.stats.by_category);
});

const handleExport = () => {
  // 创建安全的报告（不包含完整敏感内容）
  const safeReport = {
    ...props.response,
    findings: props.response.findings.map((f) => ({
      ...f,
      // 不包含原始文本，只保留预览
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
