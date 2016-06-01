package engine

import (
	"fmt"

	"github.com/egnwd/outgain/server/guest"
)

// Engineer is a mock interface for testing without real engines
type Engineer interface {
	Run()
	AddEntity(name *guest.Guest, builder builderFunc)
}

// MockEngine allows us to use an engine without acutally running a simulation
type MockEngine struct{}

// Run prints a running message to the console showing that the
// method has been called
func (m *MockEngine) Run() {
	fmt.Println("Running Engine...")
}

// AddEntity prints a message saying that the Entity has been added to the engine
func (m *MockEngine) AddEntity(_ *guest.Guest, _ builderFunc) {
	fmt.Println("Added Entity...")
}
