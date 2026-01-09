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

	result := text
	convertedFindings := make([]types.Finding, 0, len(sortedFindings))

	for _, f := range sortedFindings {
		replacement := s.getReplacement(f)
		preview := s.getPreview(f, replacement)

		// 替换文本（使用字节索引，因为 Start 和 End 是基于字节的）
		start := f.Start
		end := f.End
		if start >= 0 && end <= len(result) && start < end {
			result = result[:start] + replacement + result[end:]
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

	return result, convertedFindings
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
		// 手机号：保留前3位和后4位，如 138****0000
		prefixLen, suffixLen = 3, 4
	case detector.CategoryEmail:
		// 邮箱：保留@前3位和@后完整域名，如 san***@example.com
		atIndex := strings.Index(text, "@")
		if atIndex > 0 {
			if atIndex > 3 {
				prefixLen = 3
			} else {
				prefixLen = atIndex
			}
			// 保留@符号和完整域名部分
			suffixLen = length - atIndex
		} else {
			prefixLen, suffixLen = 3, 3
		}
	case detector.CategoryIDCard:
		// 身份证：保留前3位和后4位，如 110***XXXX
		prefixLen, suffixLen = 3, 4
	case detector.CategoryBankCard, detector.CategoryCreditCard:
		// 银行卡/信用卡：保留前4位和后4位，中间全部打码，如 6222****0000
		prefixLen, suffixLen = 4, 4
	case detector.CategoryCVV:
		// CVV：全部打码，如 ***
		return "***"
	case detector.CategoryPassport:
		// 护照号：保留首字母和后2位，如 E******00
		prefixLen, suffixLen = 1, 2
	case detector.CategoryDriverLicense:
		// 驾照号：保留前2位和后3位，如 BJ*****001
		prefixLen, suffixLen = 2, 3
	case detector.CategoryAddress:
		// 地址：保留省市区，打码详细地址，如 北京市朝阳区****
		// 简化处理：保留前6个字符，其余打码
		if length > 6 {
			prefixLen = 6
			suffixLen = 0
		} else {
			prefixLen = 2
			suffixLen = 0
		}
	case detector.CategoryGPS:
		// GPS坐标：全部打码，如 **.**, **.**
		return "**.**, **.**"
	case detector.CategoryMAC:
		// MAC地址：保留前2段，其余打码，如 00:00:**:**:**:**
		parts := strings.FieldsFunc(text, func(r rune) bool {
			return r == ':' || r == '-'
		})
		if len(parts) >= 2 {
			return parts[0] + ":" + parts[1] + ":**:**:**:**"
		}
		prefixLen, suffixLen = 2, 0
	case detector.CategoryDatabaseConn:
		// 数据库连接串：保留协议和主机，打码用户名密码，如 postgres://***:***@host:port/db
		protocolIndex := strings.Index(text, "://")
		atIndex := strings.Index(text, "@")
		if protocolIndex > 0 && atIndex > protocolIndex {
			protocolPart := text[:protocolIndex+3]
			hostPart := text[atIndex+1:]
			return protocolPart + "***:***@" + hostPart
		}
		// 如果没有@符号，只保留协议部分
		if protocolIndex > 0 {
			return text[:protocolIndex+3] + "***"
		}
		prefixLen, suffixLen = 4, 4
	case detector.CategoryName:
		// 姓名：保留首字，其余打码，如 张**
		prefixLen, suffixLen = 1, 0
	case detector.CategoryDate:
		// 日期：打码月份和日期，如 1999-**-**
		parts := strings.FieldsFunc(text, func(r rune) bool {
			return r == '-' || r == '/' || r == '.'
		})
		if len(parts) >= 3 {
			return parts[0] + "-**-**"
		}
		prefixLen, suffixLen = 4, 0
	case detector.CategoryToken:
		// Token/密钥：根据不同类型采用不同脱敏策略
		return s.maskToken(text)
	case detector.CategoryPassword:
		// 密码：全部打码，如 ********
		return strings.Repeat("*", length)
	case detector.CategoryPrivateKey:
		// 私钥：全部替换为占位符
		return "[PRIVATE_KEY_REDACTED]"
	default:
		prefixLen, suffixLen = 2, 2
	}

	if prefixLen+suffixLen >= length {
		return strings.Repeat("*", length)
	}

	prefix := text[:prefixLen]
	suffix := ""
	if suffixLen > 0 {
		suffix = text[length-suffixLen:]
	}
	masked := strings.Repeat("*", length-prefixLen-suffixLen)
	return prefix + masked + suffix
}

// maskToken 针对不同类型的Token/密钥采用精细脱敏策略
func (s *Sanitizer) maskToken(text string) string {
	length := len(text)
	if length <= 4 {
		return "****"
	}

	textUpper := strings.ToUpper(text)
	
	// AWS Access Key ID: AKIA开头，固定20字符
	if strings.HasPrefix(textUpper, "AKIA") && length == 20 {
		return "AKIA" + strings.Repeat("*", 16)
	}

	// AWS Secret Access Key: 通常40字符的base64字符串
	if length == 40 && isBase64Like(text) {
		return strings.Repeat("*", 40)
	}

	// OpenAI API Key: sk-开头
	if strings.HasPrefix(text, "sk-") || strings.HasPrefix(text, "SK-") {
		// sk-test-xxx 或 sk-xxx 格式
		parts := strings.Split(text, "-")
		if len(parts) >= 2 {
			// 保留前缀和最后4位，如 sk-****FAKE
			if length > 8 {
				return parts[0] + "-" + strings.Repeat("*", length-len(parts[0])-5) + text[length-4:]
			}
			return parts[0] + "-" + strings.Repeat("*", length-len(parts[0])-1)
		}
		return "sk-" + strings.Repeat("*", length-3)
	}

	// GitHub Token: ghp_, gho_, ghu_, ghs_, ghr_ 开头
	if strings.HasPrefix(textUpper, "GHP_") || strings.HasPrefix(textUpper, "GHO_") ||
		strings.HasPrefix(textUpper, "GHU_") || strings.HasPrefix(textUpper, "GHS_") ||
		strings.HasPrefix(textUpper, "GHR_") {
		prefix := text[:4]
		if length > 8 {
			return prefix + strings.Repeat("*", length-8) + text[length-4:]
		}
		return prefix + strings.Repeat("*", length-4)
	}

	// Stripe Token: pk_ 或 sk_ 开头
	if strings.HasPrefix(textUpper, "PK_") || strings.HasPrefix(textUpper, "SK_") {
		prefix := text[:3]
		if length > 7 {
			return prefix + strings.Repeat("*", length-7) + text[length-4:]
		}
		return prefix + strings.Repeat("*", length-3)
	}

	// JWT Token: eyJ 开头
	if strings.HasPrefix(text, "eyJ") {
		// JWT通常很长，只保留前6位和后6位
		if length > 12 {
			return text[:6] + strings.Repeat("*", length-12) + text[length-6:]
		}
		return strings.Repeat("*", length)
	}

	// 通用Token: 保留前4位和后4位
	if length > 8 {
		return text[:4] + strings.Repeat("*", length-8) + text[length-4:]
	}
	return strings.Repeat("*", length)
}

// isBase64Like 检查字符串是否像base64编码（简单检查）
func isBase64Like(s string) bool {
	base64Chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	for _, char := range s {
		if !strings.ContainsRune(base64Chars, char) {
			return false
		}
	}
	return true
}

// redact 替换为占位符
func (s *Sanitizer) redact(category detector.Category) string {
	categoryMap := map[detector.Category]string{
		detector.CategoryPhone:         "[REDACTED:PHONE]",
		detector.CategoryEmail:         "[REDACTED:EMAIL]",
		detector.CategoryIDCard:        "[REDACTED:ID_CARD]",
		detector.CategoryIP:            "[REDACTED:IP]",
		detector.CategoryDomain:        "[REDACTED:DOMAIN]",
		detector.CategoryToken:         "[REDACTED:TOKEN]",
		detector.CategoryPassword:      "[REDACTED:PASSWORD]",
		detector.CategoryPrivateKey:    "[REDACTED:PRIVATE_KEY]",
		detector.CategoryBankCard:      "[REDACTED:BANK_CARD]",
		detector.CategoryCreditCard:   "[REDACTED:CREDIT_CARD]",
		detector.CategoryCVV:           "[REDACTED:CVV]",
		detector.CategoryPassport:      "[REDACTED:PASSPORT]",
		detector.CategoryDriverLicense: "[REDACTED:DRIVER_LICENSE]",
		detector.CategoryAddress:       "[REDACTED:ADDRESS]",
		detector.CategoryGPS:           "[REDACTED:GPS]",
		detector.CategoryMAC:           "[REDACTED:MAC]",
		detector.CategoryDatabaseConn:  "[REDACTED:DATABASE_CONN]",
		detector.CategoryName:          "[REDACTED:NAME]",
		detector.CategoryDate:          "[REDACTED:DATE]",
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
		detector.CategoryPhone:         "PHONE",
		detector.CategoryEmail:         "EMAIL",
		detector.CategoryIDCard:        "ID_CARD",
		detector.CategoryIP:            "IP",
		detector.CategoryDomain:        "DOMAIN",
		detector.CategoryToken:         "TOKEN",
		detector.CategoryPassword:      "PASSWORD",
		detector.CategoryPrivateKey:    "PRIVATE_KEY",
		detector.CategoryBankCard:      "BANK_CARD",
		detector.CategoryCreditCard:    "CREDIT_CARD",
		detector.CategoryCVV:           "CVV",
		detector.CategoryPassport:      "PASSPORT",
		detector.CategoryDriverLicense: "DRIVER_LICENSE",
		detector.CategoryAddress:       "ADDRESS",
		detector.CategoryGPS:           "GPS",
		detector.CategoryMAC:           "MAC",
		detector.CategoryDatabaseConn:  "DATABASE_CONN",
		detector.CategoryName:          "NAME",
		detector.CategoryDate:          "DATE",
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
