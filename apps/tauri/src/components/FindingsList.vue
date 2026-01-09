<template>
  <div
    style="
      height: 200px;
      border-top: 1px solid #e0e0e0;
      background: #fff;
      display: flex;
      flex-direction: column;
    "
  >
    <div
      style="
        padding: 8px 16px;
        background: #f8f8f8;
        border-bottom: 1px solid #e0e0e0;
        display: flex;
        gap: 8px;
        align-items: center;
      "
    >
      <span style="font-weight: bold">命中列表 ({{ findings.length }})</span>
      <select
        :value="filterCategory"
        @change="(e) => filterCategory = (e.target as HTMLSelectElement).value"
        style="padding: 4px 8px; font-size: 12px"
      >
        <option value="all">全部类别</option>
        <option v-for="cat in categories" :key="cat" :value="cat">
          {{ categoryLabels[cat] || cat }}
        </option>
      </select>
    </div>
    <div style="flex: 1; overflow: auto; padding: 8px">
      <div
        v-if="filteredFindings.length === 0"
        style="padding: 16px; text-align: center; color: #999"
      >
        没有找到匹配项
      </div>
      <table v-else style="width: 100%; font-size: 12px">
        <thead>
          <tr style="background: #f5f5f5">
            <th style="padding: 8px; text-align: left">类别</th>
            <th style="padding: 8px; text-align: left">位置</th>
            <th style="padding: 8px; text-align: left">风险</th>
            <th style="padding: 8px; text-align: left">预览</th>
            <th style="padding: 8px; text-align: left">原因</th>
            <th style="padding: 8px; text-align: left">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(finding, i) in filteredFindings"
            :key="i"
            style="border-bottom: 1px solid #eee; cursor: pointer"
            @mouseenter="(e) => (e.currentTarget as HTMLTableRowElement).style.background = '#f9f9f9'"
            @mouseleave="(e) => (e.currentTarget as HTMLTableRowElement).style.background = 'transparent'"
          >
            <td style="padding: 8px">
              {{ categoryLabels[finding.type] || finding.type }}
            </td>
            <td style="padding: 8px">
              {{ finding.start }}-{{ finding.end }}
            </td>
            <td style="padding: 8px">
              <span
                :style="{
                  color: getRiskColor(finding.risk),
                  fontWeight: 'bold',
                }"
              >
                {{ finding.risk }}
              </span>
            </td>
            <td style="padding: 8px; font-family: monospace">
              {{ finding.replacement_preview }}
            </td>
            <td style="padding: 8px; font-size: 11px; color: #666">
              {{ finding.reason }}
            </td>
            <td style="padding: 8px">
              <button
                @click="onJump(finding.start, finding.end)"
                style="
                  padding: 4px 8px;
                  font-size: 11px;
                  background: #4a90e2;
                  color: white;
                "
              >
                跳转
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import type { Finding } from "../types";

interface Props {
  findings: Finding[];
  originalText: string;
  onJump: (start: number, end: number) => void;
}

const props = defineProps<Props>();

const filterCategory = ref<string>("all");

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

const categories = computed(() => {
  return Array.from(new Set(props.findings.map((f) => f.type)));
});

const filteredFindings = computed(() => {
  return filterCategory.value === "all"
    ? props.findings
    : props.findings.filter((f) => f.type === filterCategory.value);
});

const getRiskColor = (risk: number) => {
  if (risk >= 70) return "#f44336";
  if (risk >= 40) return "#ff9800";
  return "#2196f3";
};
</script>
