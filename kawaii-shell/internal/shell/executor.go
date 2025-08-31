package shell

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/creack/pty"
)

// CommandInfo holds cute information about commands
type CommandInfo struct {
	FriendlyName string
	Emoji        string
	Description  string
	IsDangerous  bool
}

// Shell represents the kawaii shell wrapper
type Shell struct {
	pty        *os.File
	cmd        *exec.Cmd
	output     chan string
	input      chan string
	done       chan bool
	lastOutput string
}

// Command translation map - making scary commands cute!
var CommandMap = map[string]CommandInfo{
	"ls":     {"Looking around", "📂", "Let's see what files are here!", false},
	"cd":     {"Moving", "🚶‍♀️", "Going to a new place!", false},
	"pwd":    {"Where am I?", "📍", "Showing our current location!", false},
	"mkdir":  {"Creating", "📁✨", "Making a new folder!", false},
	"rm":     {"Cleaning up", "🗑️", "Removing files (be careful!)", true},
	"cp":     {"Copying", "📋", "Making a copy of something!", false},
	"mv":     {"Moving", "📦", "Relocating files!", false},
	"cat":    {"Reading", "📖", "Let's see what's inside!", false},
	"grep":   {"Searching", "🔍", "Looking for something specific!", false},
	"find":   {"Exploring", "🗺️", "Searching everywhere!", false},
	"sudo":   {"Super powers", "💪", "Using special powers! Be careful! ✨", true},
	"git":    {"Version magic", "🪄", "Managing code history!", false},
	"npm":    {"Package magic", "📦", "Working with packages!", false},
	"python": {"Snake magic", "🐍", "Running Python code!", false},
	"node":   {"JavaScript magic", "⚡", "Running Node.js!", false},
}

// NewShell creates a new kawaii shell instance
func NewShell() (*Shell, error) {
	shell := &Shell{
		output: make(chan string, 100),
		input:  make(chan string, 10),
		done:   make(chan bool),
	}

	return shell, nil
}

// GetCommandInfo returns cute info about a command
func GetCommandInfo(command string) CommandInfo {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return CommandInfo{"Unknown", "❓", "I'm not sure what this does!", false}
	}

	baseCmd := parts[0]
	if info, exists := CommandMap[baseCmd]; exists {
		return info
	}

	return CommandInfo{"Running command", "⚡", fmt.Sprintf("Executing: %s", baseCmd), false}
}

// GetDefaultShell returns the default shell for the current OS
func GetDefaultShell() string {
	if runtime.GOOS == "windows" {
		if shell := os.Getenv("COMSPEC"); shell != "" {
			return shell
		}
		return "cmd.exe"
	}

	if shell := os.Getenv("SHELL"); shell != "" {
		return shell
	}

	return "/bin/bash"
}

// Start starts the shell session
func (s *Shell) Start() error {
	shell := GetDefaultShell()

	s.cmd = exec.Command(shell)
	s.cmd.Env = os.Environ()

	var err error
	s.pty, err = pty.Start(s.cmd)
	if err != nil {
		return fmt.Errorf("failed to start pty: %w", err)
	}

	// Start reading output from the shell
	go s.readOutput()

	// Start writing input to the shell
	go s.writeInput()

	return nil
}

// ExecuteCommand sends a command to the shell
func (s *Shell) ExecuteCommand(command string) error {
	if s.pty == nil {
		return fmt.Errorf("shell not started")
	}

	select {
	case s.input <- command + "\n":
		return nil
	default:
		return fmt.Errorf("input buffer full")
	}
}

// GetOutput returns the latest output from the shell
func (s *Shell) GetOutput() string {
	select {
	case output := <-s.output:
		s.lastOutput = output
		return output
	default:
		return s.lastOutput
	}
}

// Close closes the shell session
func (s *Shell) Close() error {
	if s.pty != nil {
		s.pty.Close()
	}
	if s.cmd != nil && s.cmd.Process != nil {
		s.cmd.Process.Kill()
	}
	close(s.done)
	return nil
}

// readOutput reads output from the PTY
func (s *Shell) readOutput() {
	scanner := bufio.NewScanner(s.pty)
	for scanner.Scan() {
		select {
		case s.output <- scanner.Text():
		case <-s.done:
			return
		}
	}
}

// writeInput writes input to the PTY
func (s *Shell) writeInput() {
	for {
		select {
		case input := <-s.input:
			if _, err := io.WriteString(s.pty, input); err != nil {
				return
			}
		case <-s.done:
			return
		}
	}
}
