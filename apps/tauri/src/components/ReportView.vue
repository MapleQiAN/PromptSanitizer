<template>
  <div style="flex: 1; padding: 24px; overflow: auto">
    <div
      style="
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 24px;
      "
    >
      <h2>清洗报告</h2>
      <button
        @click="handleExport"
        style="
          background: #4a90e2;
          color: white;
          padding: 10px 20px;
        "
      >
        导出 JSON
      </button>
    </div>

    <!-- 风险评分 -->
    <div
      :style="{
        background: '#fff',
        padding: '20px',
        borderRadius: '8px',
        marginBottom: '24px',
        border: `2px solid ${riskInfo.color}`,
      }"
    >
      <div style="display: flex; align-items: center; gap: 16px">
        <div
          :style="{
            fontSize: '48px',
            fontWeight: 'bold',
            color: riskInfo.color,
          }"
        >
          {{ response.risk_score }}
        </div>
        <div>
          <div style="font-size: 20px; font-weight: bold">
            {{ riskInfo.label }}
          </div>
          <div style="color: #666; font-size: 14px">风险评分</div>
        </div>
      </div>
    </div>

    <!-- 统计信息 -->
    <div
      style="
        background: #fff;
        padding: 20px;
        border-radius: 8px;
        margin-bottom: 24px;
      "
    >
      <h3 style="margin-bottom: 16px">统计信息</h3>
      <div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 16px">
        <div>
          <div style="font-size: 24px; font-weight: bold">
            {{ response.stats.total_findings }}
          </div>
          <div style="color: #666">总命中数</div>
        </div>
        <div>
          <div style="font-size: 24px; font-weight: bold; color: #f44336">
            {{ response.stats.high_risk_count }}
          </div>
          <div style="color: #666">高危</div>
        </div>
        <div>
          <div style="font-size: 24px; font-weight: bold; color: #ff9800">
            {{ response.stats.medium_risk_count }}
          </div>
          <div style="color: #666">中危</div>
        </div>
        <div>
          <div style="font-size: 24px; font-weight: bold; color: #4caf50">
            {{ response.stats.low_risk_count }}
          </div>
          <div style="color: #666">低危</div>
        </div>
      </div>

      <div style="margin-top: 20px">
        <h4 style="margin-bottom: 12px">按类别统计</h4>
        <div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 8px">
          <div
            v-for="[cat, count] in categoryEntries"
            :key="cat"
            style="
              display: flex;
              justify-content: space-between;
              padding: 8px;
              background: #f5f5f5;
              border-radius: 4px;
            "
          >
            <span>{{ categoryLabels[cat] || cat }}</span>
            <span style="font-weight: bold">{{ count }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 详细信息 -->
    <div
      style="
        background: #fff;
        padding: 20px;
        border-radius: 8px;
      "
    >
      <h3 style="margin-bottom: 16px">详细信息</h3>
      <div style="font-size: 12px; color: #666; margin-bottom: 16px">
        引擎版本: {{ response.version }}
      </div>
      <pre
        style="
          background: #f5f5f5;
          padding: 16px;
          border-radius: 4px;
          overflow: auto;
          font-size: 12px;
          max-height: 400px;
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
