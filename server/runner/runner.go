package runner

import (
	"bufio"
	"bytes"
	"github.com/docker/docker/pkg/reexec"
	"github.com/egnwd/outgain/server/protocol"
	"github.com/mitchellh/go-mruby"
	"io"
	"log"
	"net/rpc"
	"os"
)

func init() {
	reexec.Register("ai-runner", execRunner)
}

type Runner struct {
	mrb  *mruby.Mrb
	proc *mruby.MrbValue
}

func NewRunner(input io.Reader) (runner *Runner, err error) {
	mrb := mruby.NewMrb()
	defer func() {
		if err != nil {
			mrb.Close()
		}
	}()

	parser := mruby.NewParser(mrb)
	defer parser.Close()

	context := mruby.NewCompileContext(mrb)
	defer context.Close()

	buffer := new(bytes.Buffer)
	if _, err := buffer.ReadFrom(input); err != nil {
		return nil, err
	}

	if _, err := parser.Parse(buffer.String(), context); err != nil {
		return nil, err
	}

	proc := parser.GenerateCode()

	return &Runner{mrb, proc}, nil
}

func (runner *Runner) Tick(state protocol.WorldState, resp *struct{}) error {
	index := runner.mrb.ArenaSave()
	defer runner.mrb.ArenaRestore(index)

	if _, err := runner.mrb.Run(runner.proc, nil); err != nil {
		return err
	}

	return nil
}

func (runner *Runner) close() {
	runner.mrb.Close()
}

func execRunner() {
	runner, err := NewRunner(bufio.NewReader(os.Stdin))
	if err != nil {
		log.Fatalln(err)
	}
	defer runner.close()

	server := rpc.NewServer()
	if err = server.Register(runner); err != nil {
		log.Fatalln(err)
	}

	// Use inherited FD 3
	conn := os.NewFile(3, "")
	server.ServeConn(conn)
}
