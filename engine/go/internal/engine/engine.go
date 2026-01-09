package engine

import (
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
			detector.NewBankCardDetector(),
			detector.NewCreditCardDetector(),
			detector.NewCVVDetector(),
			detector.NewPassportDetector(),
			detector.NewDriverLicenseDetector(),
			detector.NewAddressDetector(),
			detector.NewGPSDetector(),
			detector.NewMACDetector(),
			detector.NewDatabaseConnDetector(),
			detector.NewNameDetector(),
			detector.NewDateDetector(),
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

// deduplicateFindings 去重并处理重叠
func (e *Engine) deduplicateFindings(findings []detector.Finding) []detector.Finding {
	// 先按优先级排序：风险高的优先，然后按长度（长的优先），最后按类型优先级
	// 类型优先级：风险高的优先
	typePriority := map[detector.Category]int{
		detector.CategoryEmail:         10,
		detector.CategoryDomain:        1,
		detector.CategoryToken:         8,
		detector.CategoryPassword:      7,
		detector.CategoryPrivateKey:    9,
		detector.CategoryIDCard:        9,
		detector.CategoryPhone:         6,
		detector.CategoryIP:            5,
		detector.CategoryBankCard:      9,
		detector.CategoryCreditCard:    9,
		detector.CategoryCVV:           9,
		detector.CategoryPassport:      8,
		detector.CategoryDriverLicense: 8,
		detector.CategoryAddress:       6,
		detector.CategoryGPS:           7,
		detector.CategoryMAC:           5,
		detector.CategoryDatabaseConn:  10,
		detector.CategoryName:          5,
		detector.CategoryDate:          6,
	}

	sorted := make([]detector.Finding, len(findings))
	copy(sorted, findings)
	
	// 排序：风险高的优先，然后按类型优先级，最后按长度
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			shouldSwap := false
			
			// 优先按风险排序
			if sorted[i].Risk < sorted[j].Risk {
				shouldSwap = true
			} else if sorted[i].Risk == sorted[j].Risk {
				// 风险相同时，按类型优先级
				priI := typePriority[sorted[i].Type]
				priJ := typePriority[sorted[j].Type]
				if priI < priJ {
					shouldSwap = true
				} else if priI == priJ {
					// 类型优先级相同时，按长度（长的优先）
					lenI := sorted[i].End - sorted[i].Start
					lenJ := sorted[j].End - sorted[j].Start
					if lenI < lenJ {
						shouldSwap = true
					}
				}
			}
			
			if shouldSwap {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// 去除重叠：如果两个 finding 重叠，保留优先级高的
	var unique []detector.Finding
	for _, f := range sorted {
		overlapped := false
		for _, existing := range unique {
			// 检查是否重叠：两个区间有交集
			if !(f.End <= existing.Start || f.Start >= existing.End) {
				overlapped = true
				break
			}
		}
		if !overlapped {
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
