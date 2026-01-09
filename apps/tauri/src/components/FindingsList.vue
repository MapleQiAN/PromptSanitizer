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
            <span class="badge" style="font-size: 9px; padding: 2px 6px;">
              {{ categoryLabels[finding.type] || finding.type }}
            </span>
          </td>
          <td style="color: var(--color-text-muted); font-size: 11px;">
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
          <td style="font-family: var(--font-mono); font-size: 11px; max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
            {{ finding.replacement_preview }}
          </td>
          <td style="color: var(--color-text-muted); font-size: 11px; max-width: 250px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
            {{ finding.reason }}
          </td>
          <td style="text-align: right;">
            <button
              class="btn-action"
              style="padding: 4px 12px; font-size: 10px; display: inline-flex;"
              @click.stop="onJump(finding.start, finding.end)"
            >
              Jump →
            </button>
          </td>
        </tr>
      </tbody>
    </table>
    
    <div
      v-if="findings.length === 0"
      style="
        padding: 40px;
        text-align: center;
        color: var(--color-text-muted);
        font-family: var(--font-mono);
        font-size: 11px;
        text-transform: uppercase;
        letter-spacing: 0.1em;
      "
    >
      No findings match the current filter
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
