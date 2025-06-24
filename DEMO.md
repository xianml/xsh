# xsh 使用演示

## 快速开始

### 1. 构建项目
```bash
make build
# 或者
go build -o xsh .
```

### 2. 配置 API 密钥
```bash
# 复制配置模板
cp config.example .env

# 编辑配置文件，设置至少一个 API 密钥
vim .env
```

### 3. 加载环境变量并启动
```bash
source .env
./xsh
```

## 使用示例

### 基本命令执行
```bash
xsh> ls -la
xsh> ll          # 支持 zsh aliases
xsh> la          # 支持所有用户自定义的 aliases
xsh> pwd
xsh> git status
```

### AI 辅助功能
```bash
# 使用 AI 获取命令建议
xsh> ai 找出当前目录下最大的 5 个文件

# AI 会建议类似以下命令：
🤖 AI suggests the following commands:
1. find . -type f -exec ls -lS {} + | head -5
2. du -a . | sort -nr | head -5
Execute command? (y/n/number):
```

### 更多 AI 使用场景
```bash
# 系统监控
xsh> ai 显示系统内存和 CPU 使用情况

# 文件操作
xsh> ai 递归删除所有 .DS_Store 文件

# Git 操作
xsh> ai 查看最近一周的提交历史

# 网络诊断
xsh> ai 检查网络连接和 DNS 解析

# 包管理
xsh> ai 更新系统包并清理缓存
```

### 模型切换
```bash
# 查看和切换 AI 模型
xsh> models

# 显示：
Available AI models:
1. openai (current)
2. claude
3. gemini
Select model (number): 2
Switched to model: claude
```

### 特殊命令
```bash
# 显示版本信息
./xsh --version

# 退出 xsh
xsh> exit
```

## 配置说明

### 环境变量配置
```bash
# 默认模型选择
XSH_MODEL=openai

# OpenAI
OPENAI_API_KEY=sk-your-key-here
OPENAI_MODEL=gpt-4

# Anthropic Claude
ANTHROPIC_API_KEY=sk-ant-your-key-here
ANTHROPIC_MODEL=claude-3-sonnet-20240229

# Google Gemini
GOOGLE_API_KEY=your-google-key-here
GOOGLE_MODEL=gemini-pro
```

## 颜色方案

xsh 使用不同颜色来区分不同类型的输出：

- 🔵 **蓝色** - 提示信息和 prompt
- 🟢 **绿色** - shell 命令（粗体）
- 🟡 **黄色** - AI 响应
- 🔴 **红色** - 错误信息
- 🩵 **青色** - 用户输入

## 高级功能

### 命令提取算法
xsh 使用智能算法从 AI 响应中提取 shell 命令：
- 识别以 `$` 或 `>` 开头的命令行
- 解析代码块中的 bash/sh 命令
- 匹配常见的命令模式

### 安全特性
- 所有命令在执行前都会提示用户确认
- 支持选择多个建议命令中的任意一个
- 可以查看完整的 AI 响应再决定是否执行

### zsh 兼容性
- 支持常用的内置 aliases（`ll`, `la`, `l`, `ls`, `grep` 等）
- 通过内置预处理器确保常用命令的兼容性
- 保持与原生 zsh 相同的行为和环境
- 自动处理颜色输出和常用命令别名

### 历史记录
- 命令历史保存在 `~/.xsh_history`
- 支持上下箭头浏览历史命令
- 支持 readline 的所有标准功能

## 故障排除

### 常见问题

1. **"no AI model configured" 错误**
   - 确保设置了至少一个有效的 API 密钥
   - 检查环境变量是否正确加载

2. **API 请求失败**
   - 验证 API 密钥是否有效
   - 检查网络连接
   - 确认 API 配额是否充足

3. **命令执行失败**
   - xsh 通过 zsh 执行命令，确保 zsh 可用
   - 检查命令语法是否正确

### 调试模式
```bash
# 查看详细错误信息
XSH_DEBUG=1 ./xsh
```

## 开发和贡献

### 项目结构
```
xsh/
├── main.go              # 程序入口
├── internal/
│   ├── shell/          # Shell 核心逻辑
│   ├── ai/             # AI 客户端实现
│   └── config/         # 配置管理
├── Makefile            # 构建脚本
└── README.md           # 项目文档
```

### 添加新的 AI 提供商
1. 在 `internal/ai/` 中创建新的提供商文件
2. 实现 `Provider` 接口
3. 在 `client.go` 中注册新提供商
4. 更新配置文件模板

欢迎提交 Pull Request 和 Issue！ 