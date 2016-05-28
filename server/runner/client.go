package runner

import (
	"golang.org/x/sys/unix"
	"io"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"runtime"

	"github.com/egnwd/outgain/server/protocol"
)

type RunnerClient struct {
	client *rpc.Client
	cmd    *exec.Cmd
}

// StartRunner creates a new process to execute the AI runner.
// `code` is used as the AI's source.
func StartRunner(code string) (client *RunnerClient, err error) {
	fds, err := unix.Socketpair(unix.AF_LOCAL, unix.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}

	clientFile := os.NewFile(uintptr(fds[0]), "client")
	serverFile := os.NewFile(uintptr(fds[1]), "server")
	defer func() {
		if err != nil {
			clientFile.Close()
		}
	}()
	defer serverFile.Close()

	cmd := &exec.Cmd{
		Path:       os.Args[0],
		Args:       []string{"ai-runner"},
		ExtraFiles: []*os.File{serverFile},
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	defer stdin.Close()

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	io.WriteString(stdin, code)

	client = new(RunnerClient)
	client.client = rpc.NewClient(clientFile)
	client.cmd = cmd

	runtime.SetFinalizer(client, func(client *RunnerClient) {
		client.Close()
	})

	return client, nil
}

// Tick runs a tick in the AI runner, and waits for the result
func (client *RunnerClient) Tick(state protocol.WorldState) (TickResult, error) {
	var result TickResult
	err := client.client.Call("Runner.Tick", state, &result)

	return result, err
}

// Close closes the connection to the AI runner
func (client *RunnerClient) Close() {
	if err := client.client.Close(); err != nil {
		log.Print("Error closing runner: ", err)
	}
	if err := client.cmd.Process.Kill(); err != nil {
		log.Print("Error killing runner: ", err)
	}
	go func() {
		client.cmd.Wait()
	}()
}
