<template>
  <div>
    <table class="data-table">
      <thead>
        <tr>
          <th>Type</th>
          <th>Position</th>
          <th>Risk</th>
          <th>Preview</th>
          <th>Reason</th>
          <th style="text-align: right;">Action</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(finding, i) in findings"
          :key="i"
          @click="onJump(finding.start, finding.end)"
        >
          <td>
            <span class="badge" style="font-size: 11px; padding: 6px 12px;">
              {{ categoryLabels[finding.type] || finding.type }}
            </span>
          </td>
          <td style="color: var(--color-text-muted); font-size: 12px; font-family: var(--font-mono);">
            {{ finding.start }}-{{ finding.end }}
          </td>
          <td>
            <span
              class="risk-badge"
              :class="{
                'risk-badge--high': finding.risk >= 70,
                'risk-badge--medium': finding.risk >= 40 && finding.risk < 70,
                'risk-badge--low': finding.risk < 40,
              }"
            >
              {{ finding.risk }}
            </span>
          </td>
          <td style="font-family: var(--font-mono); font-size: 12px; max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
            {{ finding.replacement_preview }}
          </td>
          <td style="color: var(--color-text-muted); font-size: 12px; max-width: 250px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
            {{ finding.reason }}
          </td>
          <td style="text-align: right;">
            <button
              class="btn-action"
              style="padding: 6px 16px; font-size: 12px; display: inline-flex; width: auto; align-items: center; gap: 6px;"
              @click.stop="onJump(finding.start, finding.end)"
            >
              <span>Jump</span>
              <span style="font-size: 14px;">→</span>
            </button>
          </td>
        </tr>
      </tbody>
    </table>
    
    <div
      v-if="findings.length === 0"
      style="
        padding: 48px;
        text-align: center;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 12px;
      "
    >
      <div style="font-size: 36px; opacity: 0.3;">🔍</div>
      <div style="color: var(--color-text-muted); font-family: var(--font-display); font-size: 14px; font-weight: 600;">
        No findings match the current filter
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Finding } from "../types";

interface Props {
  findings: Finding[];
  originalText: string;
  onJump: (start: number, end: number) => void;
}

defineProps<Props>();

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
</script>
