package detector

import (
	"testing"
)

func TestPhoneDetector(t *testing.T) {
	detector := NewPhoneDetector()

	tests := []struct {
		name     string
		text     string
		level    string
		expected int
	}{
		{"标准手机号", "我的手机号是13812345678", "standard", 1},
		{"带分隔符", "联系我：138-1234-5678", "standard", 1},
		{"宽松模式", "13812345678", "lenient", 1},
		{"严格模式", "+86-138-1234-5678", "strict", 1},
		{"无效号码", "12345678901", "standard", 0},
		{"文本中的号码", "请拨打13812345678联系", "standard", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := detector.Detect(tt.text, tt.level)
			if len(findings) != tt.expected {
				t.Errorf("expected %d findings, got %d", tt.expected, len(findings))
			}
		})
	}
}

func TestEmailDetector(t *testing.T) {
	detector := NewEmailDetector()

	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{"标准邮箱", "联系我：test@example.com", 1},
		{"多个邮箱", "a@b.com 和 c@d.org", 2},
		{"带下划线", "user_name@domain.co.uk", 1},
		{"版本号", "版本1.2.3", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := detector.Detect(tt.text, "standard")
			if len(findings) != tt.expected {
				t.Errorf("expected %d findings, got %d", tt.expected, len(findings))
			}
		})
	}
}

func TestIDCardDetector(t *testing.T) {
	detector := NewIDCardDetector()

	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{"标准身份证", "身份证号：110101199001011234", 1},
		{"末尾X", "11010119900101123X", 1},
		{"重复数字", "111111111111111111", 0},
		{"文本中", "我的身份证是110101199001011234", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := detector.Detect(tt.text, "standard")
			if len(findings) != tt.expected {
				t.Errorf("expected %d findings, got %d", tt.expected, len(findings))
			}
		})
	}
}

func TestIPDetector(t *testing.T) {
	detector := NewIPDetector()

	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{"标准IP", "服务器IP：192.168.1.1", 1},
		{"多个IP", "192.168.1.1 和 10.0.0.1", 2},
		{"无效IP", "999.999.999.999", 0},
		{"版本号", "版本1.2.3.4", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := detector.Detect(tt.text, "standard")
			if len(findings) != tt.expected {
				t.Errorf("expected %d findings, got %d", tt.expected, len(findings))
			}
		})
	}
}

func TestTokenDetector(t *testing.T) {
	detector := NewTokenDetector()

	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{"API Key", "api_key=sk-123456789012345678901234567890", 1},
		{"Bearer Token", "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", 1},
		{"JWT", "token: eyJhbGci.eyJzdWIi.IyMzM3NDk2ODk2", 1},
		{"GitHub Token", "ghp_123456789012345678901234567890123456", 1},
		{"短字符串", "key=12345", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := detector.Detect(tt.text, "standard")
			if len(findings) != tt.expected {
				t.Errorf("expected %d findings, got %d", tt.expected, len(findings))
			}
		})
	}
}

func TestPasswordDetector(t *testing.T) {
	detector := NewPasswordDetector()

	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{"密码字段", "password=mypassword123", 1},
		{"JSON格式", `{"password": "secret123"}`, 1},
		{"占位符", "password=password", 0},
		{"短密码", "pwd=12345", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := detector.Detect(tt.text, "standard")
			if len(findings) != tt.expected {
				t.Errorf("expected %d findings, got %d", tt.expected, len(findings))
			}
		})
	}
}

func TestPrivateKeyDetector(t *testing.T) {
	detector := NewPrivateKeyDetector()

	privateKey := `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA...
-----END RSA PRIVATE KEY-----`

	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{"PEM私钥", privateKey, 1},
		{"普通文本", "这不是私钥", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			findings := detector.Detect(tt.text, "standard")
			if len(findings) != tt.expected {
				t.Errorf("expected %d findings, got %d", tt.expected, len(findings))
			}
		})
	}
}
