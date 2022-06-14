package commandRunner

import (
	"bytes"
	"os"
	"os/exec"
)

type ICommandRunner interface {
	Run(arg0 string, args ...string) ([]byte, error)
	RunAsProxy(arg0 string, args ...string)        ([]byte, error) // Runs leaving stdin and stdout from parent
}

type CommandRunner struct {}
func (c CommandRunner) Run(arg0 string, args ...string) ([]byte, error) {
	bytes, err := exec.Command(arg0, args...).Output()
	return bytes, err
}
func (c CommandRunner) RunAsProxy(arg0 string, args ...string) ([]byte, error) {
	cmd := exec.Command(arg0, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	return stdout.Bytes(), err
}

// Mocked implementation for tests
type MockedCommandRunner struct {
	CallCount   int
	LastArgs    []string
	ReturnValue string
	Error       error
}
func NewMockedCommandRunner(returnValue string, error error) MockedCommandRunner {
	return MockedCommandRunner {
		CallCount: 0,
		LastArgs: []string{},
		Error: error,
		ReturnValue: returnValue,
	}
}
func (c *MockedCommandRunner) Run(arg0 string, args ...string) ([]byte, error) {
	c.CallCount++
	c.LastArgs = append([]string{arg0}, args...)
	return []byte(c.ReturnValue), c.Error
}
func (c *MockedCommandRunner) RunAsProxy(arg0 string, args ...string) ([]byte, error) {
	c.CallCount++
	c.LastArgs = append([]string{arg0}, args...)
	return []byte(c.ReturnValue), c.Error
}
