package engine

import "fmt"

// Engineer is a mock interface for testing without real engines
type Engineer interface {
	Run(EntityList)
	Kill()
	AddEntity(builder builderFunc)
	CreateEntity(builder builderFunc) Entity
	UpdateAchievements()
	UpdateLeaderboard()
}

// MockEngine allows us to use an engine without acutally running a simulation
type MockEngine struct{}

// Run prints a running message to the console showing that the
// method has been called
func (m *MockEngine) Run(_ EntityList) {
	fmt.Println("Running Engine...")
}

// Kill prints a message signalling the end of the engine
func Kill() {
	fmt.Println("Killing Engine...")
}

// AddEntity prints a message saying that the Entity has been added to the engine
func (m *MockEngine) AddEntity(_ builderFunc) {
	fmt.Println("Added Entity...")
}

// CreateEntity prints a message saying that the Entity has been added to the engine
func (m *MockEngine) CreateEntity(_ builderFunc) Entity {
	return &Creature{}
}
