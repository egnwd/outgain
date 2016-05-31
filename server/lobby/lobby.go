package lobby

import (
	"errors"
	"math/rand"

	"github.com/egnwd/outgain/server/engine"
	"github.com/egnwd/outgain/server/user"
)

const lobbySize int = 10

var lobbies = make(map[uint64]*Lobby)

// Lobby runs its own instance of an engine, and keeps track of its users
type Lobby struct {
	ID     uint64
	Engine engine.Engineer
	Users  user.List
	size   int
}

// NewLobby creates a new lobby with its own engine and list of users
func NewLobby() (lobby *Lobby) {
	e := engine.NewEngine()
	id := newID()
	lobby = &Lobby{
		ID:     id,
		Engine: e,
		Users:  user.List{},
	}

	lobbies[lobby.ID] = lobby

	return
}

// NewLobby creates a new lobby with its own engine and list of users
func NewTestLobby(e engine.Engineer, size int) (lobby *Lobby) {
	id := newID()
	lobby = &Lobby{
		ID:     id,
		Engine: e,
		Users:  user.List{},
		size:   size,
	}

	lobbies[lobby.ID] = lobby

	return
}

func newID() uint64 {
	id := uint64(rand.Uint32())
	_, ok := lobbies[id]
	for ok {
		id = uint64(rand.Uint32())
		_, ok = lobbies[id]
	}

	return id
}

// AddUser adds the specified user to the lobby, returning an error if the
// lobby is already at capacity, and running the engine if the user is
// the first to join
func (lobby *Lobby) AddUser(user *user.User) error {
	if len(lobby.Users) == lobbySize {
		return errors.New("Lobby full")
	}
	lobby.Users = append(lobby.Users, user)
	if len(lobby.Users) == 1 {
		go lobby.Engine.Run()
	}
	return nil
}

func GetLobby(id uint64) (*Lobby, bool) {
	l, ok := lobbies[id]
	return l, ok
}
