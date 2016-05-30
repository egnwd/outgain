package runner

import (
	"golang.org/x/sys/unix"
	"io"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"runtime"

	"github.com/egnwd/outgain/server/config"
	"github.com/egnwd/outgain/server/protocol"
)

type RunnerClient struct {
	client *rpc.Client
	cmd    *exec.Cmd
}

// StartRunner creates a new process to execute the AI runner.
// `code` is used as the AI's source.
func StartRunner(config *config.Config, code string) (client *RunnerClient, err error) {
	fds, err := unix.Socketpair(unix.AF_LOCAL, unix.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}

	clientFile := os.NewFile(uintptr(fds[0]), "client")
	defer func() {
		if err != nil {
			clientFile.Close()
		}
	}()

	serverFile := os.NewFile(uintptr(fds[1]), "server")
	defer serverFile.Close()

	// runner may be sandboxed and not able to access /dev/urandom
	// We open it here and pass it a file descriptor to it
	randomFile, err := os.Open("/dev/urandom")
	if err != nil {
		return nil, err
	}
	defer randomFile.Close()

	cmd := &exec.Cmd{
		ExtraFiles: []*os.File{serverFile, randomFile},
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
	}

	if config.SandboxMode != "" {
		cmd.Path = config.SandboxBin
		cmd.Args = []string{
			config.SandboxBin,
			config.SandboxMode,
			os.Args[0],
			"ai-runner",
		}
	} else {
		cmd.Path = os.Args[0]
		cmd.Args = []string{
			"ai-runner",
		}
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
	if client.client != nil {
		if err := client.client.Close(); err != nil {
			log.Print("Error closing runner: ", err)
		}
		client.client = nil
	}
	if client.cmd != nil {
		if err := client.cmd.Process.Kill(); err != nil {
			log.Print("Error killing runner: ", err)
		}
		cmd := client.cmd
		go func() {
			cmd.Wait()
		}()

		client.cmd = nil
	}
}
