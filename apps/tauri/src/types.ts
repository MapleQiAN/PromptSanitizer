export interface Config {
  mode: "annotate" | "sanitize";
  strategy: "mask" | "redact" | "pseudonym";
  level: "lenient" | "standard" | "strict";
  enabledCategories: string[];
  allowlist: string[];
}

export interface Finding {
  type: string;
  start: number;
  end: number;
  confidence: number;
  risk: number;
  replacement: string;
  replacement_preview: string;
  reason: string;
}

export interface Stats {
  total_findings: number;
  by_category: Record<string, number>;
  high_risk_count: number;
  medium_risk_count: number;
  low_risk_count: number;
}

export interface Response {
  sanitized_text: string;
  findings: Finding[];
  stats: Stats;
  risk_score: number;
  version: string;
}
