package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/term"

	"github.com/creack/pty"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/xian/xsh/internal/ai"
	"github.com/xian/xsh/internal/config"
)

type Shell struct {
	config *config.Config
	ai     *ai.Client
	colors struct {
		Prompt   *color.Color
		Command  *color.Color
		Response *color.Color
		Error    *color.Color
	}
	ptmx         *os.File
	ctx          context.Context
	cancel       context.CancelFunc
	aiHookActive atomic.Bool
}

func NewShell(cfg *config.Config) (*Shell, error) {
	ctx, cancel := context.WithCancel(context.Background())
	shell := &Shell{
		config: cfg,
		ai:     ai.New(cfg),
		ctx:    ctx,
		cancel: cancel,
	}
	shell.colors.Prompt = color.New(color.FgCyan, color.Bold)
	shell.colors.Command = color.New(color.FgGreen)
	shell.colors.Response = color.New(color.FgYellow)
	shell.colors.Error = color.New(color.FgRed)
	return shell, nil
}

func (s *Shell) Run() error {
	defer s.cancel() // Ensure all goroutines are signaled to stop on exit

	logo := `
██╗  ██╗███████╗██╗  ██╗
╚██╗██╔╝██╔════╝██║  ██║
 ╚███╔╝ ███████╗███████║
 ██╔██╗ ╚════██║██╔══██║
██╔╝ ██╗███████║██║  ██║
╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝
`
	s.colors.Command.Println(logo)
	s.colors.Command.Println("Welcome to xsh - Your AI-Enhanced Shell")
	fmt.Println()

	// Create a temporary directory for IPC
	ipcDir, err := os.MkdirTemp("", "xsh-ipc")
	if err != nil {
		return fmt.Errorf("failed to create IPC temp dir: %w", err)
	}
	defer os.RemoveAll(ipcDir)
	promptPipePath := filepath.Join(ipcDir, "prompt_pipe")
	resultPipePath := filepath.Join(ipcDir, "result_pipe")

	// Create the named pipes (FIFO)
	if err := syscall.Mkfifo(promptPipePath, 0600); err != nil {
		return fmt.Errorf("failed to create prompt pipe: %w", err)
	}
	if err := syscall.Mkfifo(resultPipePath, 0600); err != nil {
		return fmt.Errorf("failed to create result pipe: %w", err)
	}

	userShell := os.Getenv("SHELL")
	if userShell == "" {
		userShell = "/bin/zsh" // Default to zsh
	}
	shellName := filepath.Base(userShell)

	// Create a temporary directory for shell startup files
	zdotdir, err := os.MkdirTemp("", "xsh-zdotdir")
	if err != nil {
		return fmt.Errorf("failed to create zdotdir: %w", err)
	}
	defer os.RemoveAll(zdotdir)

	// Create and write the shell-specific startup script with the hook
	if err := s.createInitScript(shellName, zdotdir, promptPipePath, resultPipePath); err != nil {
		return fmt.Errorf("failed to create shell init script: %w", err)
	}

	var c *exec.Cmd
	switch shellName {
	case "zsh":
		c = exec.Command(userShell, "-l")
		c.Env = append(os.Environ(), "ZDOTDIR="+zdotdir)
	case "bash":
		bashrcPath := filepath.Join(zdotdir, ".bashrc")
		c = exec.Command(userShell, "--rcfile", bashrcPath, "-i")
	default:
		// For unsupported shells, just run them without hooks
		c = exec.Command(userShell, "-l")
	}

	s.ptmx, err = pty.Start(c)
	if err != nil {
		return fmt.Errorf("failed to start pty: %w", err)
	}
	defer s.ptmx.Close()

	// Handle window size changes
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for {
			select {
			case <-ch:
				if err := pty.InheritSize(os.Stdin, s.ptmx); err != nil {
					// This can happen if the PTY is already closed.
				}
			case <-s.ctx.Done():
				return
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial size
	defer signal.Stop(ch)
	defer close(ch)

	// Set stdin to raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to set stdin to raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// === Set up Shell Integration ===
	// Start a goroutine to listen for AI commands from the shell hook
	go s.commandServer(promptPipePath, resultPipePath, oldState)

	// === Correct Lifecycle Management ===
	// Goroutine for handling shell output. This is the primary signal for shutdown.
	errChan := make(chan error, 1)
	go func() {
		_, err := io.Copy(os.Stdout, s.ptmx)
		errChan <- err
	}()

	// Goroutine for handling user input.
	go s.handleInput()

	// Wait for the shell to exit.
	err = <-errChan
	if err != nil && err.Error() == "EOF" {
		// This is a normal exit.
		return nil
	}
	return err
}

func (s *Shell) handleInput() {
	buf := make([]byte, 1024)
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			// If the AI hook is active, pause input forwarding to avoid race conditions.
			if s.aiHookActive.Load() {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			n, err := os.Stdin.Read(buf)
			if err != nil {
				return // Can happen on exit
			}
			if _, wErr := s.ptmx.Write(buf[:n]); wErr != nil {
				return // Can happen if PTY is closed
			}
		}
	}
}

func (s *Shell) createInitScript(shellName, zdotdir, promptPipePath, resultPipePath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not get user home directory: %w", err)
	}

	var scriptPath, userRcPath, scriptContent string
	hook := ""

	switch shellName {
	case "zsh":
		scriptPath = filepath.Join(zdotdir, ".zshrc")
		userRcPath = filepath.Join(homeDir, ".zshrc")
		hook = fmt.Sprintf(`xsh_ai_widget() { local p_pipe=%q; local r_pipe=%q; local res; echo -n "$BUFFER" > "$p_pipe"; read -r res < "$r_pipe"; if [[ -n "$res" ]]; then BUFFER=$res; CURSOR=${#res}; fi; zle redisplay; }; zle -N xsh_ai_widget; bindkey '^I' xsh_ai_widget`, promptPipePath, resultPipePath)
	case "bash":
		scriptPath = filepath.Join(zdotdir, ".bashrc")
		userRcPath = filepath.Join(homeDir, ".bashrc")
		hook = `bind -x '"\t": "echo -e \"\n\x1b[31mxsh: Full AI support is only available for zsh.\x1b[0m\""'`
	default:
		return nil // No script for unsupported shells
	}

	// Safely source the user's original rc file if it exists
	scriptContent = fmt.Sprintf(`
# xsh startup script
if [ -f %q ]; then
  source %q
fi

# xsh keybinding hook
%s
`, userRcPath, userRcPath, hook)

	return os.WriteFile(scriptPath, []byte(scriptContent), 0600)
}

func (s *Shell) commandServer(promptPipePath, resultPipePath string, originalState *term.State) {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		promptPipe, err := os.OpenFile(promptPipePath, os.O_RDONLY, 0600)
		if err != nil {
			if s.ctx.Err() != nil {
				return
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}

		userInput, err := io.ReadAll(promptPipe)
		promptPipe.Close()
		if err != nil {
			if s.ctx.Err() != nil {
				return
			}
			continue
		}

		selectedCommand := s.triggerAIHook(userInput, originalState)

		resultPipe, err := os.OpenFile(resultPipePath, os.O_WRONLY, 0600)
		if err != nil {
			continue // Shell might have already aborted waiting
		}
		// Write result, even if empty (for cancellation), so the shell unblocks.
		_, _ = resultPipe.Write([]byte(selectedCommand))
		resultPipe.Close()
	}
}

func (s *Shell) triggerAIHook(bufferSnapshot []byte, originalState *term.State) string {
	s.aiHookActive.Store(true)
	defer s.aiHookActive.Store(false)

	// Correctly manage terminal state transitions for promptui
	term.Restore(int(os.Stdin.Fd()), originalState)
	defer term.MakeRaw(int(os.Stdin.Fd()))

	userInput := string(bufferSnapshot)
	if len(userInput) == 0 {
		s.handleModelSelection()
		return "" // Model selection does not return a command
	} else {
		return s.handleAIAnalysis(userInput)
	}
}

func (s *Shell) handleModelSelection() {
	modelInfos := s.ai.GetAvailableModelInfos()
	if len(modelInfos) == 0 {
		s.colors.Error.Println("No models available.")
		return
	}

	var items []string
	selectedIndex := 0
	currentModel := s.ai.GetCurrentModel()
	for i, info := range modelInfos {
		if info.DisplayName == currentModel {
			selectedIndex = i
		}
		items = append(items, fmt.Sprintf("%s (%s)", info.DisplayName, info.Provider))
	}

	prompt := promptui.Select{
		Label:     "Select AI Model",
		Items:     items,
		CursorPos: selectedIndex,
		Size:      10,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return // User cancelled
	}

	selectedInfo := modelInfos[idx]
	if err := s.ai.SwitchModelByDisplayName(selectedInfo.DisplayName, selectedInfo.Provider); err != nil {
		s.colors.Error.Printf("Failed to switch model: %v\n", err)
	} else {
		s.colors.Command.Printf("Switched to model: %s (%s)\n", selectedInfo.DisplayName, selectedInfo.Provider)
	}
}

func (s *Shell) handleAIAnalysis(userInput string) string {
	s.colors.Response.Println("\n🤖 Asking AI for:", userInput)
	response, err := s.ai.Query(userInput)
	if err != nil {
		s.colors.Error.Printf("AI error: %v\n", err)
		return ""
	}

	userMessage, suggestions := parseAIResponse(response)

	if len(suggestions) == 0 {
		s.colors.Response.Println("AI:", response) // Show raw response if parsing fails
		return ""
	}

	if userMessage != "" {
		s.colors.Prompt.Println("💡", userMessage)
	}

	items := append([]string{"[ Cancel ]"}, suggestions...)

	prompt := promptui.Select{
		Label: "Do you want to execute one of these commands?",
		Items: items,
		Size:  10,
	}

	idx, _, err := prompt.Run()
	if err != nil || idx == 0 {
		return "" // User cancelled or chose not to execute
	}

	commandToExecute := suggestions[idx-1]

	// The selected command is returned to the shell hook for execution.
	// The promptui library itself shows the final selection, so no extra printing is needed.
	return commandToExecute
}

// parseAIResponse parses the structured response from the AI.
func parseAIResponse(response string) (userMessage string, commands []string) {
	const userMsgPrefix = "USER_MESSAGE:"
	const shellCmdPrefix = "SHELL_COMMANDS:"

	shellCmdStart := strings.Index(response, shellCmdPrefix)
	if shellCmdStart == -1 {
		return "", nil // No commands found, parsing failed.
	}

	userMsgStart := strings.Index(response, userMsgPrefix)
	if userMsgStart != -1 && userMsgStart < shellCmdStart {
		userMessage = strings.TrimSpace(response[userMsgStart+len(userMsgPrefix) : shellCmdStart])
	}

	cmdBlock := strings.TrimSpace(response[shellCmdStart+len(shellCmdPrefix):])
	cmdLines := strings.Split(cmdBlock, "\n")
	for _, cmd := range cmdLines {
		if trimmed := strings.TrimSpace(cmd); trimmed != "" {
			commands = append(commands, trimmed)
		}
	}

	return userMessage, commands
}

func (s *Shell) Goodbye() {
	// Add a newline to ensure the art starts on a fresh line
	fmt.Println()
	byeArt := `
██████╗ ██╗   ██╗███████╗
██╔══██╗╚██╗ ██╔╝██╔════╝
██████╔╝ ╚████╔╝ █████╗  
██╔══██╗  ╚██╔╝  ██╔══╝  
██████╔╝   ██║   ███████╗
╚═════╝    ╚═╝   ╚══════╝
`
	s.colors.Prompt.Println(byeArt)
}
