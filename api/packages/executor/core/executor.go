package executor

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/sudzekai/web-os-api/executor/types"
	loggingAbstractions "github.com/sudzekai/web-os-api/logging/abstractions"
)

type Executor struct {
	logger loggingAbstractions.ILogger
}

func NewExecutor(factory loggingAbstractions.ILoggerFactory) *Executor {
	return &Executor{
		logger: factory.NewLogger("executor"),
	}
}

func (ex *Executor) Execute(command string, args ...string) types.CommandResult {
	cmd := exec.Command(command, args...)

	return ex.executeCommand(cmd)
}

func (ex *Executor) ExecuteInDirectory(directory string, command string, args ...string) types.CommandResult {
	cmd := exec.Command(command, args...)

	cmd.Dir = directory

	return ex.executeCommand(cmd)
}

func (ex *Executor) executeCommand(cmd *exec.Cmd) types.CommandResult {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	ex.logDebug("выполнение команды: %s %s", cmd.Path, strings.Join(cmd.Args, " "))

	if err := cmd.Start(); err != nil {
		return types.CommandResult{Error: err}
	}

	if err := cmd.Wait(); err != nil {
		return types.CommandResult{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
			Error:  err,
		}
	}

	return types.CommandResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
}

func (ex *Executor) logDebug(message string, args ...any) {
	if ex.logger != nil {
		ex.logger.LogDebug(message, args...)
	}
}
