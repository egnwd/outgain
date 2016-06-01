package engine

import (
	"errors"
)

const lobbySize int = 10

// A user, identified by unique GitHub name
type User struct {
	name      string
	resources int
}

// Slice of users
type UserList []User

// Each lobby runs its own instance of an engine, and keeps track of its users
type Lobby struct {
	engine *Engine
	users  UserList
	nextID <-chan uint64
}

// Called when a user creates a new lobby
func NewLobby() (lobby *Lobby) {
	idChannel := make(chan uint64)
	go func() {
		var id uint64 = 0
		for {
			idChannel <- id
			id++
		}
	}()

	e := NewEngine()

	lobby = &Lobby{
		engine: e,
		users:  UserList{},
		nextID: idChannel,
	}

	return
}

// Adds the specified user to the lobby, returning an error if the lobby is
// already at capacity, and running the engine if the user is the first to join
func AddUser(lobby *Lobby, user User) error {
	if len(lobby.users) == lobbySize {
		return errors.New("Lobby full")
	}
	lobby.users = append(lobby.users, user)
	if len(lobby.users) == 1 {
		lobby.engine.Run()
	}
	return nil
}
