package engine

import (
	"fmt"
	"github.com/prompt-sanitizer/engine/internal/detector"
	"github.com/prompt-sanitizer/engine/internal/sanitizer"
	"github.com/prompt-sanitizer/engine/pkg/types"
)

const Version = "0.1.0"

// Engine 清洗引擎
type Engine struct {
	detectors []detector.Detector
}

// NewEngine 创建引擎实例
func NewEngine() *Engine {
	return &Engine{
		detectors: []detector.Detector{
			detector.NewPhoneDetector(),
			detector.NewEmailDetector(),
			detector.NewIDCardDetector(),
			detector.NewIPDetector(),
			detector.NewDomainDetector(),
			detector.NewTokenDetector(),
			detector.NewPasswordDetector(),
			detector.NewPrivateKeyDetector(),
		},
	}
}

// Process 处理清洗请求
func (e *Engine) Process(req *types.Request) (*types.Response, error) {
	// 确定启用的检测器
	enabledDetectors := e.getEnabledDetectors(req.EnabledCategories)

	// 执行检测
	allFindings := make([]detector.Finding, 0)
	for _, det := range enabledDetectors {
		findings := det.Detect(req.Text, req.Level)
		allFindings = append(allFindings, findings...)
	}

	// 应用白名单过滤
	allFindings = e.applyAllowlist(allFindings, req.Text, req.Allowlist)

	// 去重（相同位置和类型）
	allFindings = e.deduplicateFindings(allFindings)

	// 如果是标注模式，不执行清洗
	var sanitizedText string
	var convertedFindings []types.Finding

	if req.Mode == "annotate" {
		sanitizedText = req.Text
		// 转换为响应格式但不替换
		convertedFindings = make([]types.Finding, len(allFindings))
		for i, f := range allFindings {
			convertedFindings[i] = types.Finding{
				Type:               string(f.Type),
				Start:              f.Start,
				End:                f.End,
				Confidence:         f.Confidence,
				Risk:               f.Risk,
				Replacement:        f.Text, // 标注模式下不替换
				ReplacementPreview: e.getPreview(f.Text, f.Risk),
				Reason:             f.Reason,
			}
		}
	} else {
		// 执行清洗
		san := sanitizer.NewSanitizer(req.Strategy, allFindings)
		sanitizedText, convertedFindings = san.Sanitize(req.Text)
	}

	// 计算统计信息
	stats := e.calculateStats(convertedFindings)

	// 计算风险评分
	riskScore := e.calculateRiskScore(convertedFindings)

	return &types.Response{
		SanitizedText: sanitizedText,
		Findings:      convertedFindings,
		Stats:         stats,
		RiskScore:     riskScore,
		Version:       Version,
	}, nil
}

// getEnabledDetectors 获取启用的检测器
func (e *Engine) getEnabledDetectors(enabledCategories []string) []detector.Detector {
	if len(enabledCategories) == 0 {
		// 如果没有指定，启用所有
		return e.detectors
	}

	categoryMap := make(map[string]bool)
	for _, cat := range enabledCategories {
		categoryMap[cat] = true
	}

	var enabled []detector.Detector
	for _, det := range e.detectors {
		if categoryMap[string(det.Category())] {
			enabled = append(enabled, det)
		}
	}
	return enabled
}

// applyAllowlist 应用白名单
func (e *Engine) applyAllowlist(findings []detector.Finding, text string, allowlist []string) []detector.Finding {
	if len(allowlist) == 0 {
		return findings
	}

	allowed := make(map[string]bool)
	for _, item := range allowlist {
		allowed[item] = true
	}

	var filtered []detector.Finding
	for _, f := range findings {
		matchedText := text[f.Start:f.End]
		if !allowed[matchedText] {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

// deduplicateFindings 去重
func (e *Engine) deduplicateFindings(findings []detector.Finding) []detector.Finding {
	seen := make(map[string]bool)
	var unique []detector.Finding

	for _, f := range findings {
		key := fmt.Sprintf("%d-%d-%s", f.Start, f.End, f.Type)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, f)
		}
	}
	return unique
}

// calculateStats 计算统计信息
func (e *Engine) calculateStats(findings []types.Finding) types.Stats {
	stats := types.Stats{
		ByCategory: make(map[string]int),
	}

	for _, f := range findings {
		stats.TotalFindings++
		stats.ByCategory[f.Type]++

		if f.Risk >= 70 {
			stats.HighRiskCount++
		} else if f.Risk >= 40 {
			stats.MediumRiskCount++
		} else {
			stats.LowRiskCount++
		}
	}

	return stats
}

// calculateRiskScore 计算风险评分 (0-100)
func (e *Engine) calculateRiskScore(findings []types.Finding) int {
	if len(findings) == 0 {
		return 0
	}

	totalRisk := 0
	for _, f := range findings {
		totalRisk += f.Risk
	}

	avgRisk := float64(totalRisk) / float64(len(findings))

	// 考虑数量因素
	countFactor := float64(len(findings)) * 2
	if countFactor > 50 {
		countFactor = 50
	}

	score := int(avgRisk + countFactor)
	if score > 100 {
		score = 100
	}

	return score
}

// getPreview 获取预览文本
func (e *Engine) getPreview(text string, risk int) string {
	if risk >= 70 {
		if len(text) > 10 {
			return text[:3] + "***" + text[len(text)-3:]
		}
		return "***"
	}
	return text
}
