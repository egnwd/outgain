package lobby

import (
	"fmt"
	"testing"

	"github.com/egnwd/outgain/server/engine"
	"github.com/stretchr/testify/assert"
)

var mockEngine = &engine.MockEngine{}

// NewTestLobby creates a new lobby with a test engine, a specific
// size and list of guests
func NewTestLobby(e engine.Engineer, size int) (lobby *Lobby) {
	id := newID()
	lobby = &Lobby{
		ID:     id,
		Engine: e,
		Guests: generalPopulation(size),
		size:   size,
	}

	lobbies[lobby.ID] = lobby

	return
}

func TestAllowUpToMaximumLimitOfLobbyUsers(t *testing.T) {
	max := 10
	lobby := NewTestLobby(mockEngine, max)
	for i := 0; i < max; i++ {
		err := lobby.AddUser(NewUser("user"))
		assert.Nil(t, err, "User should have been able to join.")
	}

	err := lobby.AddUser(NewUser("user"))
	assert.NotNil(t, err, "User should not have been able to join.")
}

func TestPopulateRemainingSpaceWithBots(t *testing.T) {
	max := 5
	bots := 3
	users := max - bots
	lobby := NewTestLobby(mockEngine, max)

	for i := 0; i < users; i++ {
		err := lobby.AddUser(NewUser("user"))
		assert.Nil(t, err, "User should have been able to join.")
	}

	var userCount int
	var botCount int

	for _, guest := range lobby.Guests.list {
		switch guest.Type {
		case userType:
			userCount++
		case botType:
			botCount++
		}
	}

	assert.Equal(t, users, userCount, fmt.Sprintf("There should be %d users", users))
	assert.Equal(t, bots, botCount, fmt.Sprintf("There should be %d bots", bots))
	assert.True(t, len(lobby.Guests.list) == max, "Lobby should always be psudeo-full")
}

func TestAddUser(t *testing.T) {
	max := 10
	lobby := NewTestLobby(mockEngine, max)

	assert.Equal(t, 0, lobby.Guests.userSize, "Should start with no users")
	err := lobby.AddUser(NewUser("user"))
	assert.True(t, err == nil, "User should have been able to join.")
	assert.Equal(t, 1, lobby.Guests.userSize, "Should insert 1 user")
}
