package runner

import (
	"golang.org/x/sys/unix"
	"io"
	"net/rpc"
	"os"
	"os/exec"

	"github.com/egnwd/outgain/server/protocol"
)

type RunnerClient struct {
	client *rpc.Client
	cmd    *exec.Cmd
}

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

	return &RunnerClient{
		client: rpc.NewClient(clientFile),
		cmd:    cmd,
	}, nil
}

func (client *RunnerClient) Tick(state protocol.WorldState) error {
	return client.client.Call("Runner.Tick", state, nil)
}

func (client *RunnerClient) Close() error {
	return client.client.Close()
}
