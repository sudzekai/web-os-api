package executor

import (
	"bytes"
	"os/exec"
)

func Execute(command string, args ...string) CommandResult {
	cmd := exec.Command(command, args...)

	return executeCommand(cmd)
}

func ExecuteInDirectory(directory string, command string, args ...string) CommandResult {
	cmd := exec.Command(command, args...)

	cmd.Dir = directory

	return executeCommand(cmd)
}

func executeCommand(cmd *exec.Cmd) CommandResult {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return CommandResult{Error: err}
	}

	if err := cmd.Wait(); err != nil {
		return CommandResult{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
			Error:  err,
		}
	}

	return CommandResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
}
