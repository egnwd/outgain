package lobby

import "github.com/egnwd/outgain/server/engine"

const lobbySize int = 10

var lobbies = make(map[uint64]*Lobby)

// Lobby runs its own instance of an engine, and keeps track of its users
type Lobby struct {
	ID     uint64
	Engine engine.Engineer
	Guests guestList
	size   int
}

// GenerateOneLobby is temporary until lobbies is fully working
// TODO: Remove once lobbies are working
func GenerateOneLobby() (lobby *Lobby) {
	for _, lobby := range lobbies {
		return lobby
	}

	return NewLobby()
}

// NewLobby creates a new lobby with its own engine and list of guests
func NewLobby() (lobby *Lobby) {
	e := engine.NewEngine()
	id := newID()
	lobby = &Lobby{
		ID:     id,
		Engine: e,
		Guests: generalPopulation(lobbySize),
		size:   lobbySize,
	}

	lobbies[lobby.ID] = lobby

	return
}

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

//This is just for testing until it's fully implemented
const baseID uint64 = 2019968050

func newID() uint64 {
	return baseID
}

func (lobby *Lobby) startEngine() {
	for _, guest := range lobby.Guests.list {
		lobby.Engine.AddEntity(guest.name, engine.RandomCreature)
	}

	go lobby.Engine.Run()
}

// GetLobby returns the Lobby with id: `id` and if it does not exist it returns
// `(nil, false)`
func GetLobby(id uint64) (*Lobby, bool) {
	l, ok := lobbies[id]
	return l, ok
}

// DestroyLobby removes looby from the global map
func DestroyLobby(lobby *Lobby) {
	lobby.Guests.list = []*guest{}
	lobby.Guests.userSize = 0
	lobby.Engine = nil
	delete(lobbies, lobby.ID)
}

// GetLobbyIDs returns an array of all the IDs in the lobbies map
func GetLobbyIDs() []uint64 {
  ids := make([]uint64, 0, len(lobbies))
  for id := range lobbies {
    ids = append(ids, id)
  }
  return ids
}
