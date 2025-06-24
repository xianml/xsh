package shell

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/xian/xsh/internal/ai"
	"github.com/xian/xsh/internal/config"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

type Shell struct {
	config *config.Config
	ai     *ai.Client
	rl     *readline.Instance
	colors struct {
		Prompt   *color.Color
		Command  *color.Color
		Response *color.Color
		Error    *color.Color
	}
	// 状态管理
	modelSelectorTriggered    bool
	commandExecutionTriggered bool
	pendingCommands           []string
	inputChannel              chan string
}

// 内置别名映射
var builtinAliases = map[string]string{
	"ll":    "ls -lh",
	"la":    "ls -la",
	"l":     "ls",
	"grep":  "grep --color=auto",
	"egrep": "egrep --color=auto",
	"fgrep": "fgrep --color=auto",
}

func NewShell(cfg *config.Config) (*Shell, error) {
	aiClient := ai.New(cfg)

	rl, err := readline.New("")
	if err != nil {
		return nil, err
	}

	shell := &Shell{
		config:       cfg,
		ai:           aiClient,
		rl:           rl,
		inputChannel: make(chan string, 1),
	}

	// 设置颜色
	shell.colors.Prompt = color.New(color.FgCyan, color.Bold)
	shell.colors.Command = color.New(color.FgGreen)
	shell.colors.Response = color.New(color.FgYellow)
	shell.colors.Error = color.New(color.FgRed)

	return shell, nil
}

func (s *Shell) Run() error {
	defer s.rl.Close()

	// 设置信号处理
	s.setupSignalHandlers()

	fmt.Println(s.colors.Prompt.Sprint("Welcome to xsh - AI Powered Shell"))
	fmt.Printf("%s Available models: %v (current: %s)\n",
		s.colors.Prompt.Sprint("🤖"),
		s.ai.GetAvailableModels(),
		s.ai.GetCurrentModel())
	fmt.Println(s.colors.Prompt.Sprint("Keyboard Shortcuts:"))
	fmt.Println(s.colors.Response.Sprint("  - Tab (empty input): Select AI model (↑↓ + Enter)"))
	fmt.Println(s.colors.Response.Sprint("  - Tab (with input): AI assistance"))
	fmt.Println(s.colors.Response.Sprint("  - 'm' or 'models': Select AI model"))
	fmt.Println(s.colors.Response.Sprint("  - Ctrl+C: Cancel current line"))
	fmt.Println(s.colors.Response.Sprint("  - 'exit': Quit xsh"))
	fmt.Println()

	// 使用自定义输入循环处理键盘事件
	return s.inputLoop()
}

// 创建能捕获键盘事件的补全器
func (s *Shell) createKeyEventCompleter() readline.AutoCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItemDynamic(func(line string) []string {
			trimmed := strings.TrimSpace(line)

			// 如果没有输入内容，Tab 键触发模型选择
			if trimmed == "" {
				// 直接在这里处理模型选择，不使用异步
				fmt.Printf("\n")
				s.handleModelSelectorWithArrows()
				return []string{} // 不显示任何补全建议
			}

			// 如果有输入内容，Tab 键触发 AI 分析
			fmt.Printf("\n🤖 Analyzing: %s\n", line)
			s.handleAIPrompt(line)
			fmt.Print(s.colors.Prompt.Sprint("xsh> "))
			return []string{} // 不显示任何补全建议
		}),
	)
}

// 使用箭头键选择模型
func (s *Shell) handleModelSelectorWithArrows() {
	fmt.Println("🔄 Getting available models...")
	modelInfos := s.ai.GetAvailableModelInfos()
	currentModel := s.ai.GetCurrentModel()

	if len(modelInfos) == 0 {
		s.colors.Error.Println("No models available. Please check your API keys.")
		fmt.Print(s.colors.Prompt.Sprint("xsh> "))
		return
	}

	// 准备选项列表
	var items []string
	var selectedIndex int
	for i, info := range modelInfos {
		label := fmt.Sprintf("%s (%s)", info.DisplayName, info.Provider)
		if info.DisplayName == currentModel {
			label += " (current)"
			selectedIndex = i
		}
		items = append(items, label)
	}

	// 创建箭头键选择器
	prompt := promptui.Select{
		Label:     "Select AI Model",
		Items:     items,
		CursorPos: selectedIndex,
		Size:      10,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   "▸ {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "✓ {{ . | green }}",
		},
	}

	idx, _, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Print(s.colors.Prompt.Sprint("xsh> "))
			return
		}
		s.colors.Error.Printf("Selection failed: %v\n", err)
		fmt.Print(s.colors.Prompt.Sprint("xsh> "))
		return
	}

	// 切换到选择的模型
	selectedInfo := modelInfos[idx]
	if err := s.ai.SwitchModelByDisplayName(selectedInfo.DisplayName, selectedInfo.Provider); err != nil {
		s.colors.Error.Printf("Failed to switch model: %v\n", err)
	} else {
		s.colors.Command.Printf("Switched to model: %s (%s)\n", selectedInfo.DisplayName, selectedInfo.Provider)
	}

	fmt.Print(s.colors.Prompt.Sprint("xsh> "))
}

func (s *Shell) inputLoop() error {
	for {
		// 使用特殊的补全器来捕获键盘事件
		s.rl.Config.AutoComplete = s.createKeyEventCompleter()
		s.rl.SetPrompt(s.colors.Prompt.Sprint("xsh> "))

		line, err := s.rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				// Ctrl+C 按下，取消当前行
				fmt.Println("^C")
				s.modelSelectorTriggered = false    // 重置状态
				s.commandExecutionTriggered = false // 重置命令执行状态
				s.rl.SetPrompt(s.colors.Prompt.Sprint("xsh> "))
				s.rl.Refresh()
				continue
			} else if err.Error() == "EOF" {
				break
			}
			return err
		}

		line = strings.TrimSpace(line)

		// 如果正在等待模型选择
		if s.modelSelectorTriggered {
			s.handleModelSelection(line)
			continue
		}

		// 如果正在等待命令执行选择
		if s.commandExecutionTriggered {
			s.handleCommandExecution(line)
			continue
		}

		// 检查是否是模型选择命令
		if line == "m" || line == "models" {
			s.showModelSelectorWithArrows()
			continue
		}

		if line == "" {
			continue
		}

		if line == "exit" {
			break
		}

		if strings.HasPrefix(line, "ai ") {
			prompt := strings.TrimPrefix(line, "ai ")
			s.handleAIPrompt(prompt)
			continue
		}

		// 执行命令
		s.executeCommand(line)
	}

	return nil
}

// 使用箭头键的模型选择器（用于 'm' 命令）
func (s *Shell) showModelSelectorWithArrows() {
	fmt.Println("🔄 Getting available models...")
	modelInfos := s.ai.GetAvailableModelInfos()
	currentModel := s.ai.GetCurrentModel()

	if len(modelInfos) == 0 {
		s.colors.Error.Println("No models available. Please check your API keys.")
		return
	}

	// 准备选项列表
	var items []string
	var selectedIndex int
	for i, info := range modelInfos {
		label := fmt.Sprintf("%s (%s)", info.DisplayName, info.Provider)
		if info.DisplayName == currentModel {
			label += " (current)"
			selectedIndex = i
		}
		items = append(items, label)
	}

	// 创建箭头键选择器
	prompt := promptui.Select{
		Label:     "Select AI Model",
		Items:     items,
		CursorPos: selectedIndex,
		Size:      10,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   "▸ {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "✓ {{ . | green }}",
		},
	}

	idx, _, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			return
		}
		s.colors.Error.Printf("Selection failed: %v\n", err)
		return
	}

	// 切换到选择的模型
	selectedInfo := modelInfos[idx]
	if err := s.ai.SwitchModelByDisplayName(selectedInfo.DisplayName, selectedInfo.Provider); err != nil {
		s.colors.Error.Printf("Failed to switch model: %v\n", err)
	} else {
		s.colors.Command.Printf("Switched to model: %s (%s)\n", selectedInfo.DisplayName, selectedInfo.Provider)
	}
	fmt.Println()
}

// 处理模型选择输入 - 现在不再使用，保留以防回退
func (s *Shell) handleModelSelection(input string) {
	s.modelSelectorTriggered = false // 重置状态

	if input == "" {
		fmt.Println() // 空行
		return
	}

	modelInfos := s.ai.GetAvailableModelInfos()
	if idx := s.parseChoice(input); idx > 0 && idx <= len(modelInfos) {
		selectedInfo := modelInfos[idx-1]
		if err := s.ai.SwitchModelByDisplayName(selectedInfo.DisplayName, selectedInfo.Provider); err != nil {
			s.colors.Error.Printf("Failed to switch model: %v\n", err)
		} else {
			s.colors.Command.Printf("Switched to model: %s (%s)\n", selectedInfo.DisplayName, selectedInfo.Provider)
		}
	} else {
		s.colors.Error.Println("Invalid choice")
	}
	fmt.Println()
}

// 处理命令执行选择输入 - 现在不再使用，保留以防回退
func (s *Shell) handleCommandExecution(input string) {
	s.commandExecutionTriggered = false // 重置状态

	if input == "" {
		fmt.Println() // 空行
		return
	}

	choice := strings.TrimSpace(input)

	if choice == "y" && len(s.pendingCommands) > 0 {
		s.executeCommand(s.pendingCommands[0])
	} else if idx := s.parseChoice(choice); idx > 0 && idx <= len(s.pendingCommands) {
		s.executeCommand(s.pendingCommands[idx-1])
	} else if choice != "n" && choice != "" {
		s.colors.Error.Println("Invalid choice")
	}

	s.pendingCommands = nil // 清空待执行命令
	fmt.Println()
}

// 异步处理模型选择器 - 现在不再使用
func (s *Shell) handleModelSelectorAsync() {
	// 这个函数保留但不再使用
}

// 保留原来的模型选择器
func (s *Shell) showModelSelector() {
	s.showModelSelectorWithArrows()
}

func (s *Shell) parseChoice(choice string) int {
	// 尝试解析为数字
	if idx, err := strconv.Atoi(choice); err == nil {
		return idx
	}
	return 0
}

func (s *Shell) setupSignalHandlers() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		s.rl.Close()
		os.Exit(0)
	}()
}

func (s *Shell) executeCommand(cmd string) {
	// 预处理别名
	parts := strings.Fields(cmd)
	if len(parts) > 0 {
		if alias, exists := builtinAliases[parts[0]]; exists {
			// 替换别名
			aliasParts := strings.Fields(alias)
			newParts := append(aliasParts, parts[1:]...)
			cmd = strings.Join(newParts, " ")
		}
	}

	fmt.Printf("%s %s\n", s.colors.Command.Sprint("$"), cmd)

	// 执行命令
	execCmd := exec.Command("zsh", "-c", cmd)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		s.colors.Error.Printf("Command failed: %v\n", err)
	}
	fmt.Println()
}

func (s *Shell) handleAIPrompt(prompt string) {
	s.colors.Response.Println("🤖 Asking AI...")

	response, err := s.ai.Query(prompt)
	if err != nil {
		s.colors.Error.Printf("AI error: %v\n", err)
		return
	}

	// 解析 AI 响应寻找建议的命令
	suggestions := s.parseAISuggestions(response)

	if len(suggestions) > 0 {
		fmt.Println("AI suggests the following commands:")
		for i, suggestion := range suggestions {
			s.colors.Command.Printf("%d. %s\n", i+1, suggestion)
		}

		// 使用箭头键选择是否执行命令
		s.handleCommandExecutionWithArrows(suggestions)
	} else {
		// 如果没有找到命令建议，就显示原始响应
		s.colors.Response.Println(response)
		fmt.Println()
	}
}

// 使用箭头键选择是否执行命令
func (s *Shell) handleCommandExecutionWithArrows(commands []string) {
	// 准备选项
	var items []string
	items = append(items, "No, don't execute any command")
	for _, cmd := range commands {
		items = append(items, fmt.Sprintf("Execute: %s", cmd))
	}

	// 创建箭头键选择器
	prompt := promptui.Select{
		Label: "Execute command?",
		Items: items,
		Size:  len(items),
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   "▸ {{ . | cyan }}",
			Inactive: "  {{ . }}",
			Selected: "✓ {{ . | green }}",
		},
	}

	idx, _, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println()
			return
		}
		s.colors.Error.Printf("Selection failed: %v\n", err)
		return
	}

	// 如果选择了执行命令（不是第一个"No"选项）
	if idx > 0 {
		s.executeCommand(commands[idx-1])
	}

	fmt.Println()
}

func (s *Shell) parseAISuggestions(response string) []string {
	var suggestions []string
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 寻找看起来像命令的行
		if strings.HasPrefix(line, "1. ") || strings.HasPrefix(line, "2. ") ||
			strings.HasPrefix(line, "3. ") || strings.HasPrefix(line, "4. ") ||
			strings.HasPrefix(line, "- ") {
			// 提取命令部分
			parts := strings.SplitN(line, " ", 2)
			if len(parts) > 1 {
				cmd := parts[1]
				// 清理命令中的注释
				if idx := strings.Index(cmd, " #"); idx != -1 {
					cmd = cmd[:idx]
				}
				if idx := strings.Index(cmd, "   #"); idx != -1 {
					cmd = cmd[:idx]
				}
				cmd = strings.TrimSpace(cmd)
				if cmd != "" {
					suggestions = append(suggestions, cmd)
				}
			}
		}
	}

	return suggestions
}
