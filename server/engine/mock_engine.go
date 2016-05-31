package engine

import "fmt"

// Engineer is a mock interface for testing without real engines
type Engineer interface {
	Run()
}

// MockEngine allows us to use an engine without acutally running a simulation
type MockEngine struct{}

// Run prints a running message to the console showing that the
// method has been called
func (m *MockEngine) Run() {
	fmt.Println("Running Engine...")
}
