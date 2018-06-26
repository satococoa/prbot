package prbot

import (
	"bytes"
	"os"
	"os/exec"
)

// ExecCommand executes given command on given path
func ExecCommand(setting *Setting, execPath string) (string, error) {
	err := os.Chdir(execPath)
	if err != nil {
		return "", err
	}
	cmdStr := setting.command
	cmd := exec.Command("sh", "-c", cmdStr)
	buffer := new(bytes.Buffer)
	cmd.Stdout = buffer
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return string(buffer.Bytes()), nil
}
