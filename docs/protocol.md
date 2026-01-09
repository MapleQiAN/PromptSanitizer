# 通信协议文档

本文档描述 Go 引擎与前端之间的 JSON 通信协议。

## 协议概述

Go sidecar 通过 stdin/stdout 进行一次性 JSON 通信：
- 前端通过 Tauri 命令调用 Go 二进制
- 请求 JSON 通过 stdin 传入
- 响应 JSON 通过 stdout 输出

## 请求格式 (Request)

```json
{
  "text": "要清洗的文本内容",
  "mode": "sanitize" | "annotate",
  "strategy": "mask" | "redact" | "pseudonym",
  "level": "lenient" | "standard" | "strict",
  "enabled_categories": ["phone", "email", ...],
  "allowlist": ["排除的字符串1", "排除的字符串2"],
  "semantic_mode": "off" | "on"
}
```

### 字段说明

- `text` (string, 必需): 要处理的文本内容
- `mode` (string, 可选): 
  - `"sanitize"`: 清洗模式，执行替换
  - `"annotate"`: 标注模式，只识别不替换
- `strategy` (string, 可选): 清洗策略
  - `"mask"`: 部分打码，保留前后缀
  - `"redact"`: 替换为占位符 `[REDACTED:TYPE]`
  - `"pseudonym"`: 一致化替换，同一实体使用相同代号
- `level` (string, 可选): 清洗强度
  - `"lenient"`: 宽松，只识别明显高风险内容
  - `"standard"`: 标准，平衡误报和漏报
  - `"strict"`: 严格，尽可能识别所有可能的敏感信息
- `enabled_categories` (string[], 可选): 启用的检测类别列表
  - 可选值: `"phone"`, `"email"`, `"id_card"`, `"ip"`, `"domain"`, `"token"`, `"password"`, `"private_key"`
  - 如果为空数组或未提供，则启用所有类别
- `allowlist` (string[], 可选): 白名单字符串列表（精确匹配）
- `semantic_mode` (string, 可选): 语义模式（预留，默认 `"off"`）

## 响应格式 (Response)

```json
{
  "sanitized_text": "清洗后的文本",
  "findings": [
    {
      "type": "phone",
      "start": 10,
      "end": 21,
      "confidence": 0.9,
      "risk": 60,
      "replacement": "[REDACTED:PHONE]",
      "replacement_preview": "138***5678",
      "reason": "检测到中国手机号格式"
    }
  ],
  "stats": {
    "total_findings": 5,
    "by_category": {
      "phone": 2,
      "email": 1,
      "token": 2
    },
    "high_risk_count": 2,
    "medium_risk_count": 2,
    "low_risk_count": 1
  },
  "risk_score": 65,
  "version": "0.1.0"
}
```

### 字段说明

- `sanitized_text` (string): 清洗后的文本（标注模式下与原文相同）
- `findings` (array): 识别到的敏感信息列表
  - `type` (string): 类别
  - `start` (int): 原文中的起始位置（字符索引）
  - `end` (int): 原文中的结束位置（字符索引）
  - `confidence` (float): 置信度 0-1
  - `risk` (int): 风险等级 0-100
  - `replacement` (string): 替换后的文本
  - `replacement_preview` (string): 用于报告的预览（掩码形式，不泄露完整内容）
  - `reason` (string): 识别原因说明
- `stats` (object): 统计信息
  - `total_findings` (int): 总命中数
  - `by_category` (object): 按类别统计
  - `high_risk_count` (int): 高危数量（risk >= 70）
  - `medium_risk_count` (int): 中危数量（40 <= risk < 70）
  - `low_risk_count` (int): 低危数量（risk < 40）
- `risk_score` (int): 整体风险评分 0-100
- `version` (string): 引擎版本号

## 错误响应

如果处理失败，Go 引擎会输出错误 JSON：

```json
{
  "error": "错误描述信息"
}
```

## 示例

### 请求示例

```json
{
  "text": "我的手机号是13812345678，邮箱是test@example.com",
  "mode": "sanitize",
  "strategy": "redact",
  "level": "standard",
  "enabled_categories": ["phone", "email"],
  "allowlist": [],
  "semantic_mode": "off"
}
```

### 响应示例

```json
{
  "sanitized_text": "我的手机号是[REDACTED:PHONE]，邮箱是[REDACTED:EMAIL]",
  "findings": [
    {
      "type": "phone",
      "start": 6,
      "end": 17,
      "confidence": 0.9,
      "risk": 60,
      "replacement": "[REDACTED:PHONE]",
      "replacement_preview": "138***5678",
      "reason": "检测到中国手机号格式"
    },
    {
      "type": "email",
      "start": 20,
      "end": 35,
      "confidence": 0.95,
      "risk": 50,
      "replacement": "[REDACTED:EMAIL]",
      "replacement_preview": "tes***com",
      "reason": "检测到邮箱地址格式"
    }
  ],
  "stats": {
    "total_findings": 2,
    "by_category": {
      "phone": 1,
      "email": 1
    },
    "high_risk_count": 0,
    "medium_risk_count": 2,
    "low_risk_count": 0
  },
  "risk_score": 55,
  "version": "0.1.0"
}
```

## 注意事项

1. **字符索引**: `start` 和 `end` 基于 UTF-8 字符索引，不是字节索引
2. **隐私保护**: `replacement_preview` 不应包含完整的敏感内容，只显示掩码预览
3. **性能**: 对于 50k 字符的文本，处理时间应 < 1 秒
4. **稳定性**: Go 引擎崩溃时，Tauri 层应捕获错误并提示用户
