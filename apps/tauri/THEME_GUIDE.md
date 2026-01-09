# 🎨 主题风格切换指南 / Theme Style Guide

## 概述 / Overview

PromptSanitizer 现在支持**两种完全不同的设计风格**，每种风格都有深浅色模式，共计 4 种主题组合。

## 🎀 风格一：Kawaii（可爱风）

### 设计理念
- **关键词**: 活泼、可爱、温暖、友好
- **灵感来源**: Nintendo Switch UI、日本可爱文化
- **适用场景**: 个人使用、轻松环境、创意工作

### 视觉特点
- **色彩**: 粉红 (#FF6B9D) + 青蓝 (#00D9FF) + 明黄 (#FFD93D)
- **形状**: 大圆角 (12-40px)
- **装饰**: emoji (🌸✨💖⭐🦋)、浮动气泡、渐变背景
- **字体**: Fredoka (展示) + Quicksand (正文)
- **动画**: 弹跳、脉冲发光、emoji 旋转
- **氛围**: 梦幻、柔和、充满生命力

### 浅色模式
- 背景: 温暖粉白色 (#FFF4F8)
- 文字: 深灰色 (#2D2D44)
- 强调: 粉红渐变

### 深色模式
- 背景: 梦幻紫黑 (#1A1625)
- 文字: 白色
- 强调: 柔和粉红

---

## 🚀 风格二：Horizon（科技风）

### 设计理念
- **关键词**: 专业、现代、科技、简洁
- **灵感来源**: [Horizon UI Dashboard](https://horizon-ui.com/horizon-tailwind-react/admin/default)
- **适用场景**: 专业环境、团队协作、正式展示

### 视觉特点
- **色彩**: 紫色渐变 (#4318FF → #7551FF) + 绿色强调 (#01B574)
- **形状**: 适中圆角 (14-30px)
- **装饰**: 无 emoji、微妙渐变网格、玻璃态效果
- **字体**: 同 Kawaii 但权重更专业
- **动画**: 平滑过渡、微妙悬停
- **氛围**: 专业、清晰、高效

### 浅色模式
- 背景: 冷色调白 (#F4F7FE)
- 文字: 深蓝灰 (#2B3674)
- 强调: 紫色渐变
- 特效: Glassmorphism (玻璃态)

### 深色模式  
- 背景: 深蓝色 (#0B1437)
- 文字: 白色
- 强调: 亮紫色 (#7551FF)
- 特效: 玻璃态边框发光

---

## 🎯 主要差异对比

| 特性 | Kawaii 可爱风 | Horizon 科技风 |
|------|--------------|---------------|
| **主色调** | 粉红 + 青蓝 | 紫色 + 绿色 |
| **圆角大小** | 大 (12-40px) | 适中 (14-30px) |
| **emoji 装饰** | ✅ 丰富 | ❌ 无 |
| **背景效果** | 彩色气泡渐变 | 玻璃态模糊 |
| **动画风格** | 弹跳、活泼 | 平滑、微妙 |
| **阴影** | 柔和彩色 | 专业灰色 |
| **按钮样式** | 圆润渐变 + 弹跳 | 渐变 + 玻璃态 |
| **Logo 动画** | 弹跳浮动 | 静态或轻微缩放 |
| **整体氛围** | 温暖友好 | 专业高效 |

---

## 🎮 如何切换主题

### 界面操作

1. **风格切换**: 点击顶部栏的 `🎀 可爱风` / `🚀 科技风` 按钮
2. **色彩模式**: 点击 `☀️ 浅色` / `🌙 深色` 按钮
3. **语言切换**: 点击 `🌏 中文` / `🌍 EN` 按钮

### 技术实现

```typescript
// 风格类型
type ThemeStyle = "kawaii" | "horizon";
const themeStyle = ref<ThemeStyle>("kawaii");

// 切换函数
const toggleStyle = () => {
  themeStyle.value = themeStyle.value === "kawaii" ? "horizon" : "kawaii";
};

// CSS 类应用
body.classList.add('style-kawaii');  // 或 'style-horizon'
body.classList.add('theme-light');   // 或 'theme-dark'
```

---

## 🎨 CSS 架构

### 变量系统

所有颜色、间距、圆角都通过 CSS 变量定义，根据风格类自动切换：

```css
/* Kawaii 浅色 */
body.style-kawaii {
  --color-primary: #FF6B9D;
  --radius-lg: 28px;
  /* ... */
}

/* Horizon 浅色 */
body.style-horizon {
  --color-primary: #4318FF;
  --radius-lg: 30px;
  /* ... */
}
```

### 风格特定样式

使用选择器覆盖特定元素：

```css
/* Kawaii: emoji 装饰 */
body.style-kawaii .panel::before {
  content: '💫';
}

/* Horizon: 移除装饰 */
body.style-horizon .panel::before {
  content: none;
}

/* Horizon: 玻璃态效果 */
body.style-horizon .app-shell__sidebar {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
}
```

---

## 🌈 设计建议

### 何时使用 Kawaii 风格？
- ✅ 个人项目展示
- ✅ 创意/设计类工作
- ✅ 轻松友好的使用场景
- ✅ 需要降低严肃感的安全工具

### 何时使用 Horizon 风格？
- ✅ 企业环境
- ✅ 团队协作展示
- ✅ 正式会议/演示
- ✅ 需要专业形象的场合

---

## 🔧 扩展与自定义

### 添加新风格

1. 在 `App.vue` 中添加风格类型：
```typescript
type ThemeStyle = "kawaii" | "horizon" | "your-style";
```

2. 在 `styles.css` 中定义变量：
```css
body.style-your-style {
  --color-primary: #yourcolor;
  /* ... */
}
```

3. 添加风格特定的覆盖样式

### 自定义现有风格

修改 CSS 变量即可调整整个风格：

```css
body.style-horizon {
  --color-primary: #yourcolor;  /* 改变主色调 */
  --radius-md: 15px;            /* 调整圆角 */
}
```

---

## 📱 响应式支持

两种风格都完全支持响应式设计：
- 网格布局自适应
- 灵活的 flexbox
- 移动端优化（如适用）

---

## ♿ 可访问性

- ✅ 高对比度文字
- ✅ 清晰的视觉层次
- ✅ 键盘导航支持
- ✅ 语义化 HTML

---

## 🎯 最佳实践

1. **保持一致性**: 在同一会话中使用同一风格
2. **根据场景选择**: 个人用可爱风，工作用科技风
3. **深浅色适配**: 根据环境光线选择合适模式
4. **性能优先**: 避免同时加载两种风格的重资源

---

**设计哲学**: 一个工具，多种面貌 —— 让用户根据心情和场景选择最舒适的界面。

Made with 💖 and 🚀
