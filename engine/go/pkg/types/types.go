package types

// Request 表示清洗请求
type Request struct {
	Text              string   `json:"text"`
	Mode              string   `json:"mode"`               // "annotate" | "sanitize"
	Strategy          string   `json:"strategy"`           // "mask" | "redact" | "pseudonym"
	Level             string   `json:"level"`              // "lenient" | "standard" | "strict"
	EnabledCategories []string `json:"enabled_categories"` // 启用的类别列表
	Allowlist         []string `json:"allowlist"`          // 白名单字符串列表
	SemanticMode      string   `json:"semantic_mode"`      // "off" | "on" (预留，默认 off)
}

// Finding 表示一个识别到的敏感信息
type Finding struct {
	Type               string  `json:"type"`                // 类别：phone, email, id_card, ip, domain, token, password, private_key
	Start              int     `json:"start"`               // 原文中的起始位置
	End                int     `json:"end"`                 // 原文中的结束位置
	Confidence         float64 `json:"confidence"`          // 置信度 0-1
	Risk               int     `json:"risk"`                // 风险等级 0-100
	Replacement        string  `json:"replacement"`         // 替换后的文本
	ReplacementPreview string  `json:"replacement_preview"` // 用于报告的预览（掩码）
	Reason             string  `json:"reason"`              // 识别原因说明
	OriginalText       string  `json:"original_text"`       // 原始文本片段（仅用于内部，不输出到 JSON）
}

// Stats 表示统计信息
type Stats struct {
	TotalFindings   int            `json:"total_findings"`
	ByCategory      map[string]int `json:"by_category"`
	HighRiskCount   int            `json:"high_risk_count"`
	MediumRiskCount int            `json:"medium_risk_count"`
	LowRiskCount    int            `json:"low_risk_count"`
}

// Response 表示清洗响应
type Response struct {
	SanitizedText string    `json:"sanitized_text"`
	Findings      []Finding `json:"findings"`
	Stats         Stats     `json:"stats"`
	RiskScore     int       `json:"risk_score"` // 0-100
	Version       string    `json:"version"`
}
