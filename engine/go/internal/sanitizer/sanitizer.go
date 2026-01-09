package sanitizer

import (
	"github.com/google/uuid"
	"github.com/prompt-sanitizer/engine/internal/detector"
	"github.com/prompt-sanitizer/engine/pkg/types"
	"sort"
	"strings"
)

// Sanitizer 清洗器
type Sanitizer struct {
	strategy     string
	findings     []detector.Finding
	pseudonymMap map[string]string // 用于一致化替换
}

// NewSanitizer 创建清洗器
func NewSanitizer(strategy string, findings []detector.Finding) *Sanitizer {
	return &Sanitizer{
		strategy:     strategy,
		findings:     findings,
		pseudonymMap: make(map[string]string),
	}
}

// Sanitize 执行清洗
func (s *Sanitizer) Sanitize(text string) (string, []types.Finding) {
	// 按位置排序，从后往前替换避免索引偏移
	sortedFindings := make([]detector.Finding, len(s.findings))
	copy(sortedFindings, s.findings)
	sort.Slice(sortedFindings, func(i, j int) bool {
		return sortedFindings[i].Start > sortedFindings[j].Start
	})

	result := []rune(text)
	convertedFindings := make([]types.Finding, 0, len(sortedFindings))

	for _, f := range sortedFindings {
		replacement := s.getReplacement(f)
		preview := s.getPreview(f, replacement)

		// 替换文本
		start := f.Start
		end := f.End
		if start < len(result) && end <= len(result) {
			newText := string(result[:start]) + replacement + string(result[end:])
			result = []rune(newText)
		}

		convertedFindings = append(convertedFindings, types.Finding{
			Type:               string(f.Type),
			Start:              f.Start,
			End:                f.End,
			Confidence:         f.Confidence,
			Risk:               f.Risk,
			Replacement:        replacement,
			ReplacementPreview: preview,
			Reason:             f.Reason,
		})
	}

	// 反转顺序以匹配原始顺序
	for i, j := 0, len(convertedFindings)-1; i < j; i, j = i+1, j-1 {
		convertedFindings[i], convertedFindings[j] = convertedFindings[j], convertedFindings[i]
	}

	return string(result), convertedFindings
}

// getReplacement 获取替换文本
func (s *Sanitizer) getReplacement(f detector.Finding) string {
	switch s.strategy {
	case "mask":
		return s.mask(f.Text, f.Type)
	case "redact":
		return s.redact(f.Type)
	case "pseudonym":
		return s.pseudonym(f.Text, f.Type)
	default:
		return s.redact(f.Type)
	}
}

// mask 部分打码
func (s *Sanitizer) mask(text string, category detector.Category) string {
	length := len(text)
	if length <= 4 {
		return "****"
	}

	// 根据类别决定保留前后缀的长度
	var prefixLen, suffixLen int
	switch category {
	case detector.CategoryPhone:
		prefixLen, suffixLen = 3, 4
	case detector.CategoryEmail:
		prefixLen = len(text[:strings.Index(text, "@")])
		if prefixLen > 3 {
			prefixLen = 3
		}
		suffixLen = len(text) - strings.Index(text, "@")
		if suffixLen > 5 {
			suffixLen = 5
		}
	case detector.CategoryIDCard:
		prefixLen, suffixLen = 3, 4
	case detector.CategoryToken, detector.CategoryPassword:
		prefixLen, suffixLen = 4, 4
	default:
		prefixLen, suffixLen = 2, 2
	}

	if prefixLen+suffixLen >= length {
		return strings.Repeat("*", length)
	}

	prefix := text[:prefixLen]
	suffix := text[length-suffixLen:]
	masked := strings.Repeat("*", length-prefixLen-suffixLen)
	return prefix + masked + suffix
}

// redact 替换为占位符
func (s *Sanitizer) redact(category detector.Category) string {
	categoryMap := map[detector.Category]string{
		detector.CategoryPhone:      "[REDACTED:PHONE]",
		detector.CategoryEmail:      "[REDACTED:EMAIL]",
		detector.CategoryIDCard:     "[REDACTED:ID_CARD]",
		detector.CategoryIP:         "[REDACTED:IP]",
		detector.CategoryDomain:     "[REDACTED:DOMAIN]",
		detector.CategoryToken:      "[REDACTED:TOKEN]",
		detector.CategoryPassword:   "[REDACTED:PASSWORD]",
		detector.CategoryPrivateKey: "[REDACTED:PRIVATE_KEY]",
	}
	if name, ok := categoryMap[category]; ok {
		return name
	}
	return "[REDACTED]"
}

// pseudonym 一致化替换
func (s *Sanitizer) pseudonym(text string, category detector.Category) string {
	key := string(category) + ":" + text
	if pseudonym, ok := s.pseudonymMap[key]; ok {
		return pseudonym
	}

	// 生成新的代号
	categoryMap := map[detector.Category]string{
		detector.CategoryPhone:      "PHONE",
		detector.CategoryEmail:      "EMAIL",
		detector.CategoryIDCard:     "ID_CARD",
		detector.CategoryIP:         "IP",
		detector.CategoryDomain:     "DOMAIN",
		detector.CategoryToken:      "TOKEN",
		detector.CategoryPassword:   "PASSWORD",
		detector.CategoryPrivateKey: "PRIVATE_KEY",
	}

	prefix := categoryMap[category]
	if prefix == "" {
		prefix = "ENTITY"
	}

	// 使用UUID生成唯一标识（简化版，实际可以用计数器）
	id := uuid.New().String()[:8]
	pseudonym := "[" + prefix + "_" + id + "]"

	s.pseudonymMap[key] = pseudonym
	return pseudonym
}

// getPreview 获取预览文本（用于报告）
func (s *Sanitizer) getPreview(f detector.Finding, replacement string) string {
	// 对于敏感内容，只显示掩码预览
	if f.Risk >= 70 {
		text := f.Text
		if len(text) > 10 {
			return text[:3] + "***" + text[len(text)-3:]
		}
		return strings.Repeat("*", len(text))
	}
	return replacement
}
