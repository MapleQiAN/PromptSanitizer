<template>
  <div class="card card--subtle" style="height: 210px; display: flex; flex-direction: column">
    <div class="card__header">
      <div style="display: flex; align-items: center; gap: 8px">
        <span class="card__title">命中列表</span>
        <span class="badge-soft">共 {{ findings.length }} 条</span>
      </div>
      <select
        :value="filterCategory"
        @change="(e) => (filterCategory = (e.target as HTMLSelectElement).value)"
        class="select-pill"
        style="min-width: 120px; font-size: 11px"
      >
        <option value="all">全部类别</option>
        <option v-for="cat in categories" :key="cat" :value="cat">
          {{ categoryLabels[cat] || cat }}
        </option>
      </select>
    </div>
    <div style="flex: 1; overflow: auto; padding: 8px 10px 10px">
      <div
        v-if="filteredFindings.length === 0"
        style="padding: 16px; text-align: center; color: #6b7280; font-size: 12px"
      >
        没有找到匹配项
      </div>
      <table v-else class="table-modern">
        <thead>
          <tr>
            <th>类别</th>
            <th>位置</th>
            <th>风险</th>
            <th>预览</th>
            <th>原因</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(finding, i) in filteredFindings"
            :key="i"
            style="cursor: pointer"
          >
            <td>
              {{ categoryLabels[finding.type] || finding.type }}
            </td>
            <td>
              {{ finding.start }}-{{ finding.end }}
            </td>
            <td>
              <span
                :style="{
                  color: getRiskColor(finding.risk),
                  fontWeight: 'bold',
                }"
              >
                {{ finding.risk }}
              </span>
            </td>
            <td style="font-family: SF Mono, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace">
              {{ finding.replacement_preview }}
            </td>
            <td style="font-size: 11px; color: #9ca3af">
              {{ finding.reason }}
            </td>
            <td>
              <button class="btn-ghost" style="font-size: 11px; padding: 4px 10px" @click="onJump(finding.start, finding.end)">
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
