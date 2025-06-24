# xsh - AI Powered Shell

一个 AI 驱动的智能 shell，包装 zsh 并提供 AI 辅助功能。

## 功能特性

- 🔧 **完全包装 zsh**：保持所有 zsh 功能的同时添加 AI 增强，支持 aliases、functions 等配置
- 🤖 **AI 辅助**：按 Tab 键触发 AI 分析和命令建议
- 🎨 **彩色输出**：不同颜色区分用户输入、shell 命令、AI 响应和提示信息
- 🔄 **多模型支持**：支持 OpenAI GPT、Anthropic Claude、Google Gemini
- ⌨️ **键盘事件**：
  - `Tab` - 触发 AI 分析当前输入
  - `Ctrl+C` - 中断操作
  - `Ctrl+D` - 退出 xsh
  - `↑/↓` - 浏览命令历史
  - `Enter` - 正常执行命令

## 安装

### 从源码构建

```bash
git clone https://github.com/xian/xsh.git
cd xsh
make build
```

### 安装到系统

```bash
make install
```

## 配置

1. 复制配置文件模板：
```bash
cp config.example .env
```

2. 编辑 `.env` 文件，配置你的 API 密钥：

```bash
# 选择默认模型
XSH_MODEL=openai

# OpenAI 配置
OPENAI_API_KEY=sk-your-openai-api-key-here
OPENAI_MODEL=gpt-4

# Anthropic Claude 配置（可选）
ANTHROPIC_API_KEY=sk-ant-your-anthropic-api-key-here
ANTHROPIC_MODEL=claude-3-sonnet-20240229

# Google Gemini 配置（可选）
GOOGLE_API_KEY=your-google-api-key-here
GOOGLE_MODEL=gemini-pro
```

3. 加载环境变量：
```bash
source .env
```

## 使用方法

### 启动 xsh

```bash
xsh
```

### 基本操作

1. **正常命令执行**：输入命令后按 `Enter`，与普通 shell 相同
   ```
   xsh> ls -la
   xsh> git status
   ```

2. **AI 辅助**：输入描述后按 `Tab`，AI 会分析并提供命令建议
   ```
   xsh> 找出当前目录下大于 100MB 的文件 [Tab]
   ```
   AI 会建议：
   ```
   🤖 Analyzing: 找出当前目录下大于 100MB 的文件
   AI suggests the following commands:
   1. find . -type f -size +100M -exec ls -lh {} +
   Execute command? (y/n/number):
   ```
   
   或使用命令方式：
   ```
   xsh> ai 找出当前目录下大于 100MB 的文件
   ```

3. **切换模型**：使用 `models` 命令在不同 AI 模型间切换
   ```
   xsh> models
   Available AI models:
   1. openai (current)
   2. claude
   Select model (number): 2
   ```

### 示例场景

- **系统管理**：
  ```
  xsh> 检查系统内存使用情况 [Tab]
  → free -h
  ```

- **文件操作**：
  ```
  xsh> 递归搜索包含 "error" 的日志文件 [Tab]
  → grep -r "error" *.log
  ```

- **Git 操作**：
  ```
  xsh> 查看最近 5 次提交的简化日志 [Tab]
  → git log --oneline -5
  ```

## 开发

### 项目结构

```
xsh/
├── main.go                 # 主程序入口
├── internal/
│   ├── shell/             # Shell 核心功能
│   │   └── shell.go
│   ├── ai/                # AI 客户端
│   │   ├── client.go      # 统一客户端接口
│   │   ├── openai.go      # OpenAI 实现
│   │   ├── anthropic.go   # Anthropic 实现
│   │   └── google.go      # Google 实现
│   └── config/            # 配置管理
│       └── config.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 构建和测试

```bash
# 格式化代码
make fmt

# 运行测试
make test

# 构建
make build

# 完整检查
make check
```

## 环境变量

| 变量名 | 描述 | 默认值 |
|--------|------|--------|
| `XSH_MODEL` | 默认使用的 AI 模型 | `openai` |
| `OPENAI_API_KEY` | OpenAI API 密钥 | - |
| `OPENAI_BASE_URL` | OpenAI API 基础 URL | `https://api.openai.com/v1` |
| `OPENAI_MODEL` | OpenAI 模型名称 | `gpt-4` |
| `ANTHROPIC_API_KEY` | Anthropic API 密钥 | - |
| `ANTHROPIC_MODEL` | Anthropic 模型名称 | `claude-3-sonnet-20240229` |
| `GOOGLE_API_KEY` | Google API 密钥 | - |
| `GOOGLE_MODEL` | Google 模型名称 | `gemini-pro` |

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
