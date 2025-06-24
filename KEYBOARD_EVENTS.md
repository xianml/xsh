# xsh 键盘事件使用指南

## 🎯 键盘事件功能

xsh 支持多种键盘事件来提供更流畅的用户体验：

### 🔑 支持的键盘事件

| 键盘操作 | 功能 | 状态 |
|----------|------|------|
| **Tab** | 触发 AI 分析当前输入 | ✅ 已实现 |
| **Ctrl+C** | 立即退出 xsh | ✅ 已实现 |
| **Ctrl+D** | 退出 xsh | ✅ 原生支持 |
| **Ctrl+O** | 模型选择菜单 | 🚧 部分实现 |
| **↑/↓** | 浏览命令历史 | ✅ 原生支持 |

### 📝 使用方法

#### 1. Tab 键 - AI 分析
```bash
xsh> find large files [Tab]
# 触发 AI 分析，获取命令建议
🤖 Analyzing: find large files
AI suggests the following commands:
1. find . -type f -size +100M -exec ls -lh {} +
2. du -a . | sort -nr | head -10
Execute command? (y/n/number):
```

#### 2. 备用命令方式
如果 Tab 键不起作用，可以使用命令：
```bash
xsh> ai find large files
# 等同于上面的 Tab 键功能
```

#### 3. 模型切换
```bash
xsh> models
Available AI models:
1. openai (current)
2. claude  
3. gemini
Select model (number): 2
Switched to model: claude
```

## 🛠️ 技术实现

### Tab 键事件处理
```go
// 创建能捕获键盘事件的补全器
func (s *Shell) createKeyEventCompleter() readline.AutoCompleter {
    return readline.NewPrefixCompleter(
        readline.PcItemDynamic(func(line string) []string {
            // 这里捕获 Tab 键事件
            if strings.TrimSpace(line) != "" {
                // 立即处理 AI 请求
                fmt.Printf("\n🤖 Analyzing: %s\n", line)
                s.handleAIPrompt(line)
                // 重新显示提示符
                fmt.Print(s.colors.Prompt.Sprint("xsh> "))
            }
            return []string{} // 不显示任何补全建议
        }),
    )
}
```

### Ctrl+C 中断处理
```go
func (s *Shell) inputLoop() error {
    for {
        line, err := s.rl.Readline()
        if err != nil {
            if err == readline.ErrInterrupt {
                // Ctrl+C 按下，直接退出
                fmt.Println("\nGoodbye!")
                os.Exit(0)
            }
            // ... 其他错误处理
        }
        // ... 命令处理
    }
}
```

## 🎮 互动示例

### 场景 1：系统管理
```bash
xsh> check system memory [Tab]
🤖 AI suggests:
1. free -h
2. top -l 1 | grep PhysMem
3. vm_stat
Execute command? (y/1/2/3): 1
$ free -h
```

### 场景 2：文件操作
```bash
xsh> remove all .DS_Store files [Tab]
🤖 AI suggests:
1. find . -name ".DS_Store" -delete
2. find . -name ".DS_Store" -exec rm {} \;
Execute command? (y/1/2): 1
$ find . -name ".DS_Store" -delete
```

### 场景 3：Git 操作
```bash
xsh> show recent commits [Tab]
🤖 AI suggests:
1. git log --oneline -10
2. git log --graph --oneline -5
Execute command? (y/1/2): 2
$ git log --graph --oneline -5
```

## 🔧 故障排除

### Tab 键不响应
1. **检查终端兼容性**：确保使用支持的终端
2. **使用备用方式**：`ai <your question>`
3. **检查输入**：确保有输入内容才按 Tab

### 键盘快捷键冲突
- **Ctrl+C**：正常中断，不会直接退出
- **Ctrl+D**：EOF 信号，正常退出
- **其他快捷键**：使用命令方式作为备选

## 📈 性能优化

- **事件捕获**：轻量级，不影响正常输入
- **AI 调用**：异步处理，不阻塞界面
- **历史记录**：保持 readline 原生功能

## 🔮 未来计划

- [ ] Ctrl+O 键绑定（模型切换）
- [ ] Ctrl+R 键绑定（AI 搜索历史）
- [ ] 自定义键绑定配置
- [ ] 更多快捷键支持

## 💡 最佳实践

1. **描述清晰**：给 AI 清晰的任务描述
2. **确认执行**：仔细检查 AI 建议的命令
3. **备用方案**：熟悉命令方式作为备选
4. **模型选择**：根据任务选择合适的 AI 模型

---

*使用 `./xsh` 启动并体验这些键盘事件功能！* 