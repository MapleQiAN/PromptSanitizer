package detector

import (
	"regexp"
	"strings"
)

// Category 表示检测类别
type Category string

const (
	CategoryPhone      Category = "phone"
	CategoryEmail      Category = "email"
	CategoryIDCard     Category = "id_card"
	CategoryIP         Category = "ip"
	CategoryDomain     Category = "domain"
	CategoryToken      Category = "token"
	CategoryPassword   Category = "password"
	CategoryPrivateKey Category = "private_key"
)

// Finding 表示检测结果
type Finding struct {
	Type       Category
	Start      int
	End        int
	Text       string
	Confidence float64
	Risk       int
	Reason     string
}

// Detector 检测器接口
type Detector interface {
	Detect(text string, level string) []Finding
	Category() Category
}

// BaseDetector 基础检测器
type BaseDetector struct {
	category Category
}

func (d *BaseDetector) Category() Category {
	return d.category
}

// PhoneDetector 手机号检测器
type PhoneDetector struct {
	BaseDetector
}

func NewPhoneDetector() *PhoneDetector {
	return &PhoneDetector{BaseDetector{category: CategoryPhone}}
}

func (d *PhoneDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// 中国手机号：1开头，11位数字
	// 宽松：严格匹配 1[3-9]\d{9}
	// 标准：允许空格/横线分隔
	// 严格：更宽松的模式
	var pattern *regexp.Regexp
	switch level {
	case "lenient":
		pattern = regexp.MustCompile(`1[3-9]\d{9}`)
	case "strict":
		pattern = regexp.MustCompile(`(?:\+86[- ]?)?1[3-9]\d[- ]?\d{4}[- ]?\d{4}`)
	default: // standard
		pattern = regexp.MustCompile(`1[3-9]\d[- ]?\d{4}[- ]?\d{4}|1[3-9]\d{9}`)
	}

	matches := pattern.FindAllStringIndex(text, -1)
	for _, match := range matches {
		matchedText := text[match[0]:match[1]]
		// 清理分隔符后验证长度
		cleaned := regexp.MustCompile(`[- ]`).ReplaceAllString(matchedText, "")
		if len(cleaned) == 11 {
			findings = append(findings, Finding{
				Type:       CategoryPhone,
				Start:      match[0],
				End:        match[1],
				Text:       matchedText,
				Confidence: 0.9,
				Risk:       60,
				Reason:     "检测到中国手机号格式",
			})
		}
	}
	return findings
}

// EmailDetector 邮箱检测器
type EmailDetector struct {
	BaseDetector
}

func NewEmailDetector() *EmailDetector {
	return &EmailDetector{BaseDetector{category: CategoryEmail}}
}

func (d *EmailDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// 标准邮箱格式
	pattern := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	matches := pattern.FindAllStringIndex(text, -1)
	for _, match := range matches {
		matchedText := text[match[0]:match[1]]
		// 排除明显不是邮箱的（如版本号 1.2.3）
		if !strings.Contains(matchedText, "@") {
			continue
		}
		findings = append(findings, Finding{
			Type:       CategoryEmail,
			Start:      match[0],
			End:        match[1],
			Text:       matchedText,
			Confidence: 0.95,
			Risk:       50,
			Reason:     "检测到邮箱地址格式",
		})
	}
	return findings
}

// IDCardDetector 身份证号检测器
type IDCardDetector struct {
	BaseDetector
}

func NewIDCardDetector() *IDCardDetector {
	return &IDCardDetector{BaseDetector{category: CategoryIDCard}}
}

func (d *IDCardDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// 18位身份证号：前17位数字 + 最后一位数字或X
	pattern := regexp.MustCompile(`\b\d{17}[\dXx]\b`)
	matches := pattern.FindAllStringIndex(text, -1)
	for _, match := range matches {
		matchedText := text[match[0]:match[1]]
		// 简单校验：不能全是相同数字
		if isRepeatingDigits(matchedText[:17]) {
			continue
		}
		findings = append(findings, Finding{
			Type:       CategoryIDCard,
			Start:      match[0],
			End:        match[1],
			Text:       matchedText,
			Confidence: 0.85,
			Risk:       90,
			Reason:     "检测到18位身份证号格式",
		})
	}
	return findings
}

// IPDetector IP地址检测器
type IPDetector struct {
	BaseDetector
}

func NewIPDetector() *IPDetector {
	return &IPDetector{BaseDetector{category: CategoryIP}}
}

func (d *IPDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// IPv4 地址
	pattern := regexp.MustCompile(`\b(?:(?:25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d\d?)\b`)
	matches := pattern.FindAllStringIndex(text, -1)
	for _, match := range matches {
		matchedText := text[match[0]:match[1]]
		// 排除版本号（如 192.168.1.1.0）
		if strings.Count(matchedText, ".") == 3 {
			findings = append(findings, Finding{
				Type:       CategoryIP,
				Start:      match[0],
				End:        match[1],
				Text:       matchedText,
				Confidence: 0.9,
				Risk:       40,
				Reason:     "检测到IPv4地址格式",
			})
		}
	}
	return findings
}

// DomainDetector 域名检测器
type DomainDetector struct {
	BaseDetector
}

func NewDomainDetector() *DomainDetector {
	return &DomainDetector{BaseDetector{category: CategoryDomain}}
}

func (d *DomainDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// URL 和域名
	urlPattern := regexp.MustCompile(`https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(?:/[^\s]*)?`)
	domainPattern := regexp.MustCompile(`\b[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}\b`)

	// 检测 URL
	urlMatches := urlPattern.FindAllStringIndex(text, -1)
	for _, match := range urlMatches {
		matchedText := text[match[0]:match[1]]
		findings = append(findings, Finding{
			Type:       CategoryDomain,
			Start:      match[0],
			End:        match[1],
			Text:       matchedText,
			Confidence: 0.95,
			Risk:       30,
			Reason:     "检测到URL格式",
		})
	}

	// 检测纯域名（排除已匹配的URL）
	domainMatches := domainPattern.FindAllStringIndex(text, -1)
	for _, match := range domainMatches {
		// 检查是否已被URL匹配覆盖
		overlapped := false
		for _, f := range findings {
			if match[0] >= f.Start && match[1] <= f.End {
				overlapped = true
				break
			}
		}
		if !overlapped {
			matchedText := text[match[0]:match[1]]
			// 排除常见非敏感域名（可根据需要扩展）
			if !isCommonDomain(matchedText) {
				findings = append(findings, Finding{
					Type:       CategoryDomain,
					Start:      match[0],
					End:        match[1],
					Text:       matchedText,
					Confidence: 0.7,
					Risk:       20,
					Reason:     "检测到域名格式",
				})
			}
		}
	}
	return findings
}

// TokenDetector Token/Key检测器
type TokenDetector struct {
	BaseDetector
}

func NewTokenDetector() *TokenDetector {
	return &TokenDetector{BaseDetector{category: CategoryToken}}
}

func (d *TokenDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// API Key 常见前缀
	apiKeyPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:api[_-]?key|apikey)[\s:=]+([a-zA-Z0-9_\-]{20,})`),
		regexp.MustCompile(`(?i)(?:bearer|token)[\s:]+([a-zA-Z0-9_\-\.]{20,})`),
		regexp.MustCompile(`(?i)sk-[a-zA-Z0-9]{20,}`),  // OpenAI style
		regexp.MustCompile(`(?i)pk_[a-zA-Z0-9]{20,}`),  // Stripe style
		regexp.MustCompile(`(?i)ghp_[a-zA-Z0-9]{20,}`), // GitHub token
	}

	// JWT
	jwtPattern := regexp.MustCompile(`eyJ[a-zA-Z0-9_-]+\.eyJ[a-zA-Z0-9_-]+\.[a-zA-Z0-9_-]+`)

	// Cookie
	cookiePattern := regexp.MustCompile(`(?i)(?:cookie|sessionid)[\s:=]+([a-zA-Z0-9_\-=]{20,})`)

	allPatterns := []struct {
		pattern *regexp.Regexp
		reason  string
		risk    int
	}{
		{jwtPattern, "检测到JWT Token格式", 80},
		{cookiePattern, "检测到Cookie/Session格式", 70},
	}

	for _, p := range apiKeyPatterns {
		allPatterns = append(allPatterns, struct {
			pattern *regexp.Regexp
			reason  string
			risk    int
		}{p, "检测到API Key格式", 85})
	}

	for _, item := range allPatterns {
		matches := item.pattern.FindAllStringSubmatchIndex(text, -1)
		for _, match := range matches {
			start := match[0]
			end := match[1]
			if len(match) > 2 && match[2] >= 0 {
				// 有捕获组，使用捕获组范围
				start = match[2]
				end = match[3]
			}
			if end-start < 20 {
				continue
			}
			matchedText := text[start:end]
			findings = append(findings, Finding{
				Type:       CategoryToken,
				Start:      start,
				End:        end,
				Text:       matchedText,
				Confidence: 0.9,
				Risk:       item.risk,
				Reason:     item.reason,
			})
		}
	}
	return findings
}

// PasswordDetector 密码字段检测器
type PasswordDetector struct {
	BaseDetector
}

func NewPasswordDetector() *PasswordDetector {
	return &PasswordDetector{BaseDetector{category: CategoryPassword}}
}

func (d *PasswordDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// 密码字段模式
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:password|pwd|passwd)[\s:=]+([^\s\n]{6,})`),
		regexp.MustCompile(`(?i)"password"\s*:\s*"([^"]{6,})"`),
		regexp.MustCompile(`(?i)'password'\s*:\s*'([^']{6,})'`),
	}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatchIndex(text, -1)
		for _, match := range matches {
			if len(match) >= 4 && match[2] >= 0 {
				start := match[2]
				end := match[3]
				matchedText := text[start:end]
				// 排除明显的占位符
				if !isPlaceholder(matchedText) {
					findings = append(findings, Finding{
						Type:       CategoryPassword,
						Start:      start,
						End:        end,
						Text:       matchedText,
						Confidence: 0.85,
						Risk:       95,
						Reason:     "检测到密码字段",
					})
				}
			}
		}
	}
	return findings
}

// PrivateKeyDetector 私钥检测器
type PrivateKeyDetector struct {
	BaseDetector
}

func NewPrivateKeyDetector() *PrivateKeyDetector {
	return &PrivateKeyDetector{BaseDetector{category: CategoryPrivateKey}}
}

func (d *PrivateKeyDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// PEM 格式私钥
	pemPattern := regexp.MustCompile(`-----BEGIN [A-Z ]+ PRIVATE KEY-----\s*[A-Za-z0-9+/=\s]+\s*-----END [A-Z ]+ PRIVATE KEY-----`)
	matches := pemPattern.FindAllStringIndex(text, -1)
	for _, match := range matches {
		matchedText := text[match[0]:match[1]]
		findings = append(findings, Finding{
			Type:       CategoryPrivateKey,
			Start:      match[0],
			End:        match[1],
			Text:       matchedText,
			Confidence: 0.99,
			Risk:       100,
			Reason:     "检测到PEM格式私钥",
		})
	}
	return findings
}

// 辅助函数

func isRepeatingDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	first := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != first {
			return false
		}
	}
	return true
}

func isCommonDomain(domain string) bool {
	common := []string{"localhost", "example.com", "test.com", "example.org"}
	domainLower := strings.ToLower(domain)
	for _, c := range common {
		if strings.Contains(domainLower, c) {
			return true
		}
	}
	return false
}

func isPlaceholder(text string) bool {
	placeholders := []string{"password", "******", "***", "xxx", "placeholder", "your_password"}
	textLower := strings.ToLower(text)
	for _, p := range placeholders {
		if textLower == p {
			return true
		}
	}
	return false
}
