package runner

import (
	"errors"
	"golang.org/x/sys/unix"
	"log"
	"net/rpc"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/ugorji/go/codec"

	"github.com/egnwd/outgain/server/config"
	"github.com/egnwd/outgain/server/protocol"
)

const runnerTimeout = 50 * time.Millisecond

type RunnerClient struct {
	client *rpc.Client
	cmd    *exec.Cmd
}

// StartRunner creates a new process to execute the AI runner.
// `code` is used as the AI's source.
func StartRunner(config *config.Config) (client *RunnerClient, err error) {
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

	cmd := exec.Command(config.RunnerBin)
	cmd.ExtraFiles = []*os.File{serverFile, randomFile}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	client = new(RunnerClient)
	var mh codec.MsgpackHandle
	codec := codec.MsgpackSpecRpc.ClientCodec(clientFile, &mh)
	client.client = rpc.NewClientWithCodec(codec)
	client.cmd = cmd

	runtime.SetFinalizer(client, func(client *RunnerClient) {
		client.Close()
	})

	return client, nil
}

func (client *RunnerClient) call(method string, args interface{}, reply interface{}) error {
	call := client.client.Go(method, args, reply, nil)
	select {
	case <-call.Done:
		return call.Error
	case <-time.After(runnerTimeout):
		return errors.New("AI timed out")
	}
}

func (client *RunnerClient) Load(source string) error {
	return client.call("Runner.Load", source, nil)
}

// Tick runs a tick in the AI runner, and waits for the result
func (client *RunnerClient) Tick(player protocol.Entity, state protocol.WorldState) (protocol.TickResult, error) {
	request := protocol.TickRequest{
		WorldState: state,
		Player:     player,
	}

	var result protocol.TickResult
	err := client.call("Runner.Tick", request, &result)

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
