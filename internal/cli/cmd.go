package cli

import (
	"io"
	"os"
	"os/exec"
)

type CmdResult struct {
	Result string
	Error  error
}

func ExecuteCmd(cmd string, args ...string) (result string, err error) {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	return string(out), err
}

func ExecuteCmdAsync(cmd string, args ...string) <-chan CmdResult {
	result := make(chan CmdResult)
	go func() {
		out, err := exec.Command(cmd, args...).CombinedOutput()
		result <- CmdResult{string(out), err}
	}()
	return result
}

func ExecuteCmdToStdout(name string, args ...string) error {
	cmd := exec.Command(name, args...) // replace with your command

	// Get the pipe for command's standard output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	// Start the command
	if err = cmd.Start(); err != nil {
		return err
	}

	// Copy the data written to the Pipe to os.Stdout
	if _, err = io.Copy(os.Stdout, stdout); err != nil {
		return err
	}

	// Wait for the command to finish
	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func ExecuteCommandInTerminal(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
