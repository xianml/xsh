# xsh - AI Powered Shell

一个AI驱动的智能shell，采用**Shell Hook架构**，在你的原生shell上添加AI功能。

## 🎯 革命性Shell Hook架构

### 💡 设计理念
xsh采用**真正的Shell Hook**设计：
- 🔧 **启动原生shell进程**：`$SHELL -i` 完全加载用户环境
- 🎣 **智能Hook拦截**：只在Tab键时介入AI功能
- 🔄 **透明代理模式**：所有其他操作直接传递给原生shell
- ✨ **零配置兼容**：无需解析.zshrc/.bashrc，天然支持所有别名和函数

### 🚀 核心优势

**100% Shell兼容性**：
- ✅ 所有用户自定义别名（ll, la, grep, 等）
- ✅ Shell函数和脚本
- ✅ 环境变量和PATH配置
- ✅ 命令历史和自动补全
- ✅ Shell提示符样式
- ✅ zsh/bash所有特性

**无侵入AI增强**：
- 🎯 Tab键（空输入）：AI模型选择
- 🎯 Tab键（有输入）：AI命令分析
- 🔄 其他所有操作：原生shell行为
- ⚡ 零延迟，按需AI

## 🏗️ 技术架构

### Shell Hook实现
```go
// 启动原生shell进程
s.shellCmd = exec.Command(userShell, "-i")
s.shellStdin, _ = s.shellCmd.StdinPipe()
s.shellStdout, _ = s.shellCmd.StdoutPipe()

// Hook拦截器
func (s *Shell) createHookCompleter() readline.AutoCompleter {
    return readline.PcItemDynamic(func(line string) []string {
        if strings.TrimSpace(line) == "" {
            s.handleModelSelection()  // AI模型选择
        } else {
            s.handleAIAnalysis(line)  // AI命令分析
        }
        return []string{}
    })
}

// 透明传递
s.shellStdin.Write([]byte(line + "\n"))
```

### AI集成流程
1. **监听输入**：readline捕获用户键盘输入
2. **Hook判断**：Tab键触发AI，其他直接传递
3. **AI处理**：调用AI API分析命令需求
4. **结果展示**：箭头键选择执行建议
5. **命令执行**：选择的命令传递给原生shell

## 🎮 使用体验

### 原生Shell体验
```bash
# 完全正常的shell操作，零学习成本
xsh> ll                          # 你的别名
xsh> cd ~/projects && pwd        # 内置命令
xsh> git status | grep modified  # 管道操作
xsh> source ~/.zshrc            # 配置重载
xsh> history | tail -10         # 历史命令
```

### AI增强体验
```bash
# 需要帮助时按Tab键
xsh> find files larger than 100MB [Tab]
🤖 AI suggests:
1. find . -type f -size +100M -exec ls -lh {} +
2. du -a . | awk '$1 > 100*1024' | sort -nr
Execute command? ▸ Execute: find . -type f -size +100M -exec ls -lh {} +

# 空输入选择模型
xsh> [Tab]
Select AI Model:
▸ claude-3-sonnet-20240229 (Anthropic) (current)
  gpt-3.5-turbo (OpenAI)
  gemini-pro (Google)
```

## 🔧 完全兼容性验证

### 别名支持测试
```bash
# 所有用户别名都完全支持
xsh> ll          # ls -lh
xsh> la          # ls -la  
xsh> grep        # grep --color=auto
xsh> myalias     # 任何自定义别名
```

### 高级功能测试
```bash
# Shell函数
xsh> myfunc() { echo "Hello $1"; }
xsh> myfunc World

# 环境变量
xsh> export MY_VAR=test
xsh> echo $MY_VAR

# 复杂命令
xsh> ps aux | grep chrome | awk '{print $2}' | xargs kill
```

## ⚡ 性能特点

- **启动速度**：与原生shell相同
- **执行延迟**：零延迟（直接传递）
- **内存使用**：轻量级Hook层
- **AI调用**：按需触发，不影响日常使用

## 🎯 与传统方案对比

| 特性 | 传统Shell包装 | xsh Shell Hook |
|-----|-------------|---------------|
| 别名支持 | 需要解析配置 | ✅ 天然支持 |
| 函数支持 | 不支持 | ✅ 完全支持 |
| 环境继承 | 部分支持 | ✅ 100%继承 |
| 配置兼容 | 需要维护 | ✅ 零配置 |
| 性能开销 | 命令解析 | ✅ 零开销 |
| 用户体验 | 学习成本 | ✅ 透明使用 |

## 🎯 核心特性

- 🔧 **完全Shell兼容**：透明代理模式，100%兼容你的zsh/bash环境
  - 所有别名、函数、环境变量完全保持
  - 内置命令、管道、重定向正常工作
  - 继承用户的.zshrc/.bashrc配置
  - **不触发AI时，就是原生shell**
- 🤖 **智能AI辅助**：按Tab键获取命令建议和帮助
- 🎨 **美观界面**：彩色输出和箭头键选择
- 🔄 **多AI模型**：支持OpenAI、Anthropic、Google
- ⌨️ **直观操作**：Tab键智能行为，空输入选模型，有输入问AI

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

## 🚀 完全Shell兼容

### 核心设计理念
xsh采用**透明shell代理**设计：
- 每个命令通过 `$SHELL -i -c` 执行
- 完全继承你的shell环境和配置
- 保持所有原生shell特性
- **只在明确触发AI功能时才介入**

### 兼容性验证
```bash
# 所有这些都完全正常工作：
xsh> ll                    # 你的ll别名
xsh> cd ~ && pwd          # 内置命令和链式操作
xsh> echo $PATH           # 环境变量
xsh> ls | grep *.go       # 管道和通配符
xsh> export VAR=test && echo $VAR  # 变量导出
xsh> which zsh            # 内置which命令
xsh> source ~/.zshrc      # source命令
```

### AI功能触发
- **普通命令**：直接执行，完全兼容原shell
- **AI模式**：
  - Tab键（空输入）：选择AI模型
  - Tab键（有输入）：AI分析和建议
  - `ai <问题>`：直接AI对话
  - `m`：模型选择

## 🎯 使用场景

### 日常Shell使用
```bash
# 完全正常的shell操作
xsh> cd /path/to/project
xsh> git status
xsh> npm install
xsh> ll
```

### AI辅助场景
```bash
# 需要帮助时按Tab
xsh> find large files [Tab]
🤖 AI suggests:
1. find . -type f -size +100M -exec ls -lh {} +
2. du -a . | sort -nr | head -10
Execute command? ▸ Execute: find . -type f -size +100M -exec ls -lh {} +
```

### 模型管理
```bash
# 快速切换AI模型
xsh> [Tab]    # 空输入时按Tab
Select AI Model:
▸ gpt-3.5-turbo (OpenAI)
  claude-3-sonnet-20240229 (Anthropic) (current)
  gemini-pro (Google)
```

## ⚡ 性能和兼容性

### 性能特点
- **零延迟**：普通命令执行速度与原shell相同
- **按需AI**：只在需要时调用AI服务
- **内存友好**：不常驻后台进程

### 兼容性测试
经过以下场景测试：
- ✅ zsh/bash别名和函数
- ✅ 环境变量和PATH
- ✅ 管道和重定向
- ✅ 后台任务和进程控制
- ✅ shell内置命令
- ✅ 复杂命令链
- ✅ 脚本执行

## 🔧 故障排除

### 常见问题

1. **命令不可用/别名不工作**
   ```bash
   # 检查shell配置是否正确加载
   xsh> echo $SHELL
   xsh> source ~/.zshrc    # 重新加载配置
   ```

2. **环境变量问题**
   ```bash
   # xsh完全继承环境变量
   xsh> env | grep YOUR_VAR
   ```

3. **AI功能不响应**
   ```bash
   # 检查API密钥设置
   xsh> echo $ANTHROPIC_API_KEY
   ```

### 调试模式
```bash
# 查看详细执行信息
XSH_DEBUG=1 ./xsh
```

## 📈 技术架构

### 透明代理实现
```go
// 核心执行逻辑
func (s *Shell) executeInNativeShell(cmd string) {
    userShell := os.Getenv("SHELL")
    execCmd := exec.Command(userShell, "-i", "-c", cmd)
    execCmd.Stdin = os.Stdin
    execCmd.Stdout = os.Stdout  
    execCmd.Stderr = os.Stderr
    execCmd.Env = os.Environ()  // 完全继承环境
    execCmd.Run()
}
```

### AI集成设计
- **事件驱动**：Tab键触发AI分析
- **智能解析**：从AI响应中提取命令
- **用户确认**：所有AI建议都需要用户确认执行

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
