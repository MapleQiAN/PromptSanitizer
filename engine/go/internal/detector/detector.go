package detector

import (
	"regexp"
	"strings"
)

// Category 表示检测类别
type Category string

const (
	CategoryPhone         Category = "phone"
	CategoryEmail         Category = "email"
	CategoryIDCard        Category = "id_card"
	CategoryIP            Category = "ip"
	CategoryDomain        Category = "domain"
	CategoryToken         Category = "token"
	CategoryPassword      Category = "password"
	CategoryPrivateKey    Category = "private_key"
	CategoryBankCard      Category = "bank_card"
	CategoryCreditCard    Category = "credit_card"
	CategoryCVV           Category = "cvv"
	CategoryPassport      Category = "passport"
	CategoryDriverLicense Category = "driver_license"
	CategoryAddress       Category = "address"
	CategoryGPS           Category = "gps"
	CategoryMAC           Category = "mac"
	CategoryDatabaseConn  Category = "database_conn"
	CategoryName          Category = "name"
	CategoryDate          Category = "date"
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

// BankCardDetector 银行卡号检测器
type BankCardDetector struct {
	BaseDetector
}

func NewBankCardDetector() *BankCardDetector {
	return &BankCardDetector{BaseDetector{category: CategoryBankCard}}
}

func (d *BankCardDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// 中国银行卡号：16-19位数字，可能包含空格或横线
	// 常见格式：6222 0000 0000 0000 或 6222-0000-0000-0000
	pattern := regexp.MustCompile(`\b(?:\d{4}[- ]?){3,4}\d{1,4}\b`)
	matches := pattern.FindAllStringIndex(text, -1)
	for _, match := range matches {
		matchedText := text[match[0]:match[1]]
		cleaned := regexp.MustCompile(`[- ]`).ReplaceAllString(matchedText, "")
		// 银行卡号通常是16-19位
		if len(cleaned) >= 16 && len(cleaned) <= 19 {
			// 排除明显不是银行卡号的（如IP地址、版本号等）
			if !isIPLike(cleaned) {
				findings = append(findings, Finding{
					Type:       CategoryBankCard,
					Start:      match[0],
					End:        match[1],
					Text:       matchedText,
					Confidence: 0.8,
					Risk:       90,
					Reason:     "检测到银行卡号格式",
				})
			}
		}
	}
	return findings
}

// CreditCardDetector 信用卡号检测器
type CreditCardDetector struct {
	BaseDetector
}

func NewCreditCardDetector() *CreditCardDetector {
	return &CreditCardDetector{BaseDetector{category: CategoryCreditCard}}
}

func (d *CreditCardDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// 信用卡号：13-19位数字，可能包含空格或横线
	// Visa: 4开头，MasterCard: 5开头，Amex: 3开头
	pattern := regexp.MustCompile(`\b(?:3[47]\d|4\d|5[1-5]\d)[- ]?\d{4}[- ]?\d{4}[- ]?\d{4,7}\b`)
	matches := pattern.FindAllStringIndex(text, -1)
	for _, match := range matches {
		matchedText := text[match[0]:match[1]]
		cleaned := regexp.MustCompile(`[- ]`).ReplaceAllString(matchedText, "")
		if len(cleaned) >= 13 && len(cleaned) <= 19 {
			findings = append(findings, Finding{
				Type:       CategoryCreditCard,
				Start:      match[0],
				End:        match[1],
				Text:       matchedText,
				Confidence: 0.85,
				Risk:       95,
				Reason:     "检测到信用卡号格式",
			})
		}
	}
	return findings
}

// CVVDetector CVV检测器
type CVVDetector struct {
	BaseDetector
}

func NewCVVDetector() *CVVDetector {
	return &CVVDetector{BaseDetector{category: CategoryCVV}}
}

func (d *CVVDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// CVV通常是3-4位数字，出现在"CVV"、"CVC"、"安全码"等关键词附近
	pattern := regexp.MustCompile(`(?i)(?:cvv|cvc|安全码|验证码)[\s:：]+(\d{3,4})`)
	matches := pattern.FindAllStringSubmatchIndex(text, -1)
	for _, match := range matches {
		if len(match) >= 4 && match[2] >= 0 {
			start := match[2]
			end := match[3]
			matchedText := text[start:end]
			findings = append(findings, Finding{
				Type:       CategoryCVV,
				Start:      start,
				End:        end,
				Text:       matchedText,
				Confidence: 0.9,
				Risk:       95,
				Reason:     "检测到CVV安全码",
			})
		}
	}
	return findings
}

// PassportDetector 护照号检测器
type PassportDetector struct {
	BaseDetector
}

func NewPassportDetector() *PassportDetector {
	return &PassportDetector{BaseDetector{category: CategoryPassport}}
}

func (d *PassportDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// 中国护照：E/G/D/S开头 + 8-9位数字或字母
	// 国际护照：字母开头 + 数字和字母组合
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:护照号|passport)[\s:：]+([EGDS]\d{7,9})`),
		regexp.MustCompile(`\b[EGDS]\d{7,9}\b`),
		regexp.MustCompile(`\b[A-Z]\d{6,9}\b`),
	}
	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatchIndex(text, -1)
		for _, match := range matches {
			start := match[0]
			end := match[1]
			if len(match) > 2 && match[2] >= 0 {
				start = match[2]
				end = match[3]
			}
			matchedText := text[start:end]
			// 排除明显不是护照号的（如版本号、产品编号等）
			if !isVersionNumber(matchedText) {
				findings = append(findings, Finding{
					Type:       CategoryPassport,
					Start:      start,
					End:        end,
					Text:       matchedText,
					Confidence: 0.75,
					Risk:       85,
					Reason:     "检测到护照号格式",
				})
			}
		}
	}
	return findings
}

// DriverLicenseDetector 驾照号检测器
type DriverLicenseDetector struct {
	BaseDetector
}

func NewDriverLicenseDetector() *DriverLicenseDetector {
	return &DriverLicenseDetector{BaseDetector{category: CategoryDriverLicense}}
}

func (d *DriverLicenseDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// 中国驾照号：18位，格式类似身份证
	// 也可能包含地区代码和编号，如：BJDL-FAKE-2024-00001
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(?:驾照号|驾驶证号|driver[_\s]?license)[\s:：]+([A-Z]{2,4}[-_][A-Z0-9]+[-_]\d{4}[-_]\d{5,})`),
		regexp.MustCompile(`\b[A-Z]{2,4}[-_][A-Z0-9]+[-_]\d{4}[-_]\d{5,}\b`),
		regexp.MustCompile(`\b\d{18}\b`), // 18位数字格式
	}
	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatchIndex(text, -1)
		for _, match := range matches {
			start := match[0]
			end := match[1]
			if len(match) > 2 && match[2] >= 0 {
				start = match[2]
				end = match[3]
			}
			matchedText := text[start:end]
			if !isRepeatingDigits(matchedText) {
				findings = append(findings, Finding{
					Type:       CategoryDriverLicense,
					Start:      start,
					End:        end,
					Text:       matchedText,
					Confidence: 0.8,
					Risk:       80,
					Reason:     "检测到驾照号格式",
				})
			}
		}
	}
	return findings
}

// AddressDetector 地址检测器
type AddressDetector struct {
	BaseDetector
}

func NewAddressDetector() *AddressDetector {
	return &AddressDetector{BaseDetector{category: CategoryAddress}}
}

func (d *AddressDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// 中国地址模式：省市区/县 + 街道/路 + 门牌号 + 室/号
	// 例如：北京市朝阳区某某路88号 6号楼1203室
	pattern := regexp.MustCompile(`(?:[^，。\n]{0,10}?(?:省|市|区|县|镇|乡|街道|路|街|道|号|室|楼|层|座)[^，。\n]{0,20}?){2,}`)
	matches := pattern.FindAllStringIndex(text, -1)
	for _, match := range matches {
		matchedText := text[match[0]:match[1]]
		// 排除太短或太长的匹配（可能是误报）
		if len(matchedText) >= 10 && len(matchedText) <= 100 {
			// 必须包含地址关键词
			hasKeyword := regexp.MustCompile(`(?:省|市|区|县|镇|乡|街道|路|街|道|号|室|楼|层|座)`).MatchString(matchedText)
			if hasKeyword {
				findings = append(findings, Finding{
					Type:       CategoryAddress,
					Start:      match[0],
					End:        match[1],
					Text:       matchedText,
					Confidence: 0.7,
					Risk:       60,
					Reason:     "检测到地址信息",
				})
			}
		}
	}
	return findings
}

// GPSDetector GPS坐标检测器
type GPSDetector struct {
	BaseDetector
}

func NewGPSDetector() *GPSDetector {
	return &GPSDetector{BaseDetector{category: CategoryGPS}}
}

func (d *GPSDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// GPS坐标格式：纬度,经度 或 (纬度, 经度)
	// 例如：39.9042, 116.4074 或 (39.9042, 116.4074)
	pattern := regexp.MustCompile(`(?:GPS|坐标|经纬度|latitude|longitude)[\s:：]+(?:\(?)(\d{1,3}\.\d{1,6})[\s,，]+(\d{1,3}\.\d{1,6})(?:\)?)`)
	matches := pattern.FindAllStringSubmatchIndex(text, -1)
	for _, match := range matches {
		if len(match) >= 6 && match[2] >= 0 && match[4] >= 0 {
			lat := text[match[2]:match[3]]
			lon := text[match[4]:match[5]]
			// 验证坐标范围（粗略）
			if isValidGPS(lat, lon) {
				start := match[0]
				end := match[1]
				matchedText := text[start:end]
				findings = append(findings, Finding{
					Type:       CategoryGPS,
					Start:      start,
					End:        end,
					Text:       matchedText,
					Confidence: 0.85,
					Risk:       70,
					Reason:     "检测到GPS坐标",
				})
			}
		}
	}
	// 也检测简单的坐标格式（没有关键词前缀）
	simplePattern := regexp.MustCompile(`\b(\d{1,3}\.\d{4,6})[\s,，]+(\d{1,3}\.\d{4,6})\b`)
	simpleMatches := simplePattern.FindAllStringSubmatchIndex(text, -1)
	for _, match := range simpleMatches {
		if len(match) >= 6 && match[2] >= 0 && match[4] >= 0 {
			lat := text[match[2]:match[3]]
			lon := text[match[4]:match[5]]
			if isValidGPS(lat, lon) {
				start := match[0]
				end := match[1]
				matchedText := text[start:end]
				findings = append(findings, Finding{
					Type:       CategoryGPS,
					Start:      start,
					End:        end,
					Text:       matchedText,
					Confidence: 0.6,
					Risk:       70,
					Reason:     "检测到GPS坐标格式",
				})
			}
		}
	}
	return findings
}

// MACDetector MAC地址检测器
type MACDetector struct {
	BaseDetector
}

func NewMACDetector() *MACDetector {
	return &MACDetector{BaseDetector{category: CategoryMAC}}
}

func (d *MACDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// MAC地址格式：XX:XX:XX:XX:XX:XX 或 XX-XX-XX-XX-XX-XX
	pattern := regexp.MustCompile(`\b(?:[0-9A-Fa-f]{2}[:-]){5}[0-9A-Fa-f]{2}\b`)
	matches := pattern.FindAllStringIndex(text, -1)
	for _, match := range matches {
		matchedText := text[match[0]:match[1]]
		// 排除全0或全F的MAC地址（可能是占位符）
		if !isPlaceholderMAC(matchedText) {
			findings = append(findings, Finding{
				Type:       CategoryMAC,
				Start:      match[0],
				End:        match[1],
				Text:       matchedText,
				Confidence: 0.9,
				Risk:       50,
				Reason:     "检测到MAC地址格式",
			})
		}
	}
	return findings
}

// DatabaseConnDetector 数据库连接串检测器
type DatabaseConnDetector struct {
	BaseDetector
}

func NewDatabaseConnDetector() *DatabaseConnDetector {
	return &DatabaseConnDetector{BaseDetector{category: CategoryDatabaseConn}}
}

func (d *DatabaseConnDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// 数据库连接串：postgres://, mysql://, mongodb://, redis://等
	pattern := regexp.MustCompile(`(?i)(?:postgres|mysql|mongodb|redis|sqlserver|oracle)://[^\s\)]+`)
	matches := pattern.FindAllStringIndex(text, -1)
	for _, match := range matches {
		matchedText := text[match[0]:match[1]]
		// 连接串通常包含用户名、密码、主机、端口、数据库名
		if strings.Contains(matchedText, "@") || strings.Contains(matchedText, ":") {
			findings = append(findings, Finding{
				Type:       CategoryDatabaseConn,
				Start:      match[0],
				End:        match[1],
				Text:       matchedText,
				Confidence: 0.95,
				Risk:       100,
				Reason:     "检测到数据库连接串",
			})
		}
	}
	return findings
}

// NameDetector 姓名检测器（中文姓名）
type NameDetector struct {
	BaseDetector
}

func NewNameDetector() *NameDetector {
	return &NameDetector{BaseDetector{category: CategoryName}}
}

func (d *NameDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// 中文姓名：2-4个汉字，出现在"姓名"、"名字"等关键词附近
	// 或者在紧急联系人、家庭成员等上下文中
	namePattern := regexp.MustCompile(`(?i)(?:姓名|名字|联系人|成员|配偶|父亲|母亲|同学)[\s:：]+([\u4e00-\u9fa5]{2,4})`)
	matches := namePattern.FindAllStringSubmatchIndex(text, -1)
	for _, match := range matches {
		if len(match) >= 4 && match[2] >= 0 {
			start := match[2]
			end := match[3]
			matchedText := text[start:end]
			// 排除常见非姓名词汇
			if !isCommonWord(matchedText) {
				findings = append(findings, Finding{
					Type:       CategoryName,
					Start:      start,
					End:        end,
					Text:       matchedText,
					Confidence: 0.8,
					Risk:       50,
					Reason:     "检测到姓名信息",
				})
			}
		}
	}
	return findings
}

// DateDetector 日期检测器（用于出生日期等敏感日期）
type DateDetector struct {
	BaseDetector
}

func NewDateDetector() *DateDetector {
	return &DateDetector{BaseDetector{category: CategoryDate}}
}

func (d *DateDetector) Detect(text string, level string) []Finding {
	var findings []Finding
	// 日期格式：YYYY-MM-DD, YYYY/MM/DD, YYYY.MM.DD
	// 出现在"出生日期"、"生日"等关键词附近
	datePattern := regexp.MustCompile(`(?i)(?:出生日期|生日|birthday|date of birth)[\s:：]+(\d{4}[-/.]\d{1,2}[-/.]\d{1,2})`)
	matches := datePattern.FindAllStringSubmatchIndex(text, -1)
	for _, match := range matches {
		if len(match) >= 4 && match[2] >= 0 {
			start := match[2]
			end := match[3]
			matchedText := text[start:end]
			// 验证日期合理性（粗略）
			if isValidDate(matchedText) {
				findings = append(findings, Finding{
					Type:       CategoryDate,
					Start:      start,
					End:        end,
					Text:       matchedText,
					Confidence: 0.85,
					Risk:       60,
					Reason:     "检测到出生日期",
				})
			}
		}
	}
	return findings
}

// 辅助函数

func isIPLike(s string) bool {
	// 检查是否像IP地址（4段数字，每段0-255）
	ipPattern := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	return ipPattern.MatchString(s)
}

func isVersionNumber(s string) bool {
	// 检查是否像版本号（如 E1.2.3）
	versionPattern := regexp.MustCompile(`^[A-Z]\d+\.\d+`)
	return versionPattern.MatchString(s)
}

func isValidGPS(lat, lon string) bool {
	// 粗略验证GPS坐标范围
	// 纬度：-90 到 90，经度：-180 到 180
	// 中国大致范围：纬度 18-54，经度 73-135
	return true // 简化实现，实际可以添加更严格的验证
}

func isPlaceholderMAC(mac string) bool {
	// 检查是否是占位符MAC地址
	cleaned := regexp.MustCompile(`[:-]`).ReplaceAllString(mac, "")
	return cleaned == "000000000000" || cleaned == "FFFFFFFFFFFF"
}

func isCommonWord(word string) bool {
	// 常见非姓名词汇
	common := []string{"某某", "测试", "示例", "用户", "管理员", "系统"}
	for _, c := range common {
		if word == c {
			return true
		}
	}
	return false
}

func isValidDate(dateStr string) bool {
	// 粗略验证日期格式
	// 实际可以添加更严格的日期验证
	return regexp.MustCompile(`^\d{4}[-/.]\d{1,2}[-/.]\d{1,2}$`).MatchString(dateStr)
}

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
