package runner

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/rpc"
	"os"

	"github.com/egnwd/outgain/server/protocol"

	"github.com/docker/docker/pkg/reexec"
	"github.com/mitchellh/go-mruby"
)

func init() {
	reexec.Register("ai-runner", func() {
		// Execute the runner, reading the AI source from stdin, and server
		// on inherited FD 3
		ExecRunner(bufio.NewReader(os.Stdin), os.NewFile(3, "connection"))
	})
}

// ExecRunner starts a new AI runner, reading the AI source from
// `input`, and exposes a `net/rpc` interface on `conn`
func ExecRunner(input io.Reader, conn io.ReadWriteCloser) {
	runner, err := NewRunner(input)
	if err != nil {
		log.Fatalln(err)
	}
	defer runner.close()

	server := rpc.NewServer()
	if err = server.Register(runner); err != nil {
		log.Fatalln(err)
	}

	server.ServeConn(conn)

	fmt.Println("Runner done")
}

// Runner wraps an mruby instance loaded with an AI's source.
// It expose a `Tick` method over a `net/rpc` interface for the engine to call.
type Runner struct {
	mrb    *mruby.Mrb
	result *TickResult
}

// NewRunner creates a new AI runner, loading the AI's source from `input`.
func NewRunner(input io.Reader) (runner *Runner, err error) {
	runner = new(Runner)
	runner.mrb = mruby.NewMrb()
	defer func() {
		if err != nil {
			runner.mrb.Close()
		}
	}()

	index := runner.mrb.ArenaSave()
	defer runner.mrb.ArenaRestore(index)

	runner.mrb.ObjectClass().DefineMethod("move", runner.moveMethod(), mruby.ArgsReq(2))

	if err := runner.seedRNG(); err != nil {
		return nil, fmt.Errorf("seed: %v", err)
	}

	bytes, err := ioutil.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("read: %v", err)
	}

	if _, err := runner.mrb.LoadString(string(bytes)); err != nil {
		return nil, fmt.Errorf("load: %v", err)
	}

	return runner, nil
}

// TickResult represents the desired action to be taken by the AI
type TickResult struct {
	Dx, Dy float64
}

// Tick executes the AI to determine the desired action
func (runner *Runner) Tick(state protocol.WorldState, resp *TickResult) error {
	if resp != nil {
		resp.Dx = 0
		resp.Dy = 0
	}

	runner.result = resp

	index := runner.mrb.ArenaSave()
	defer runner.mrb.ArenaRestore(index)

	if _, err := runner.mrb.TopSelf().Call("run"); err != nil {
		return fmt.Errorf("run: %v", err)
	}

	return nil
}

func (runner *Runner) close() {
	runner.mrb.Close()
}

func valueToFloat(v *mruby.MrbValue) (float64, error) {
	switch v.Type() {
	case mruby.TypeFixnum:
		return float64(v.Fixnum()), nil
	case mruby.TypeFloat:
		return v.Float(), nil
	default:
		return 0, fmt.Errorf("Expected number")
	}
}

// seedRNG initializes mruby's RNG's seed, using the OS' random source
func (runner *Runner) seedRNG() error {
	bigSeed, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return err
	}
	seed := int(bigSeed.Int64())

	_, err = runner.mrb.TopSelf().Call("srand", mruby.Int(seed))
	return err
}

// moveMethod is exposed as `move` to the AI.
// It takes in two arguments, the desired speed on the x and y axes.
// It saves the speed inside the runner, which later gets return at the end of Tick
func (runner *Runner) moveMethod() mruby.Func {
	return func(mrb *mruby.Mrb, self *mruby.MrbValue) (mruby.Value, mruby.Value) {
		args := mrb.GetArgs()

		dx, err := valueToFloat(args[0])
		if err != nil {
			log.Fatalln(err)
		}
		dy, err := valueToFloat(args[1])
		if err != nil {
			log.Fatalln(err)
		}

		if runner.result != nil {
			runner.result.Dx = dx
			runner.result.Dy = dy
		}
		return nil, nil
	}
}
