package lobby

import (
	"errors"
	"fmt"

	"sync"

	"github.com/egnwd/outgain/server/config"
	"github.com/egnwd/outgain/server/engine"
	"github.com/egnwd/outgain/server/guest"
)

const lobbySize int = 10

var lobbies = make(map[uint64]*Lobby)

// Lobby runs its own instance of an engine, and keeps track of its users
type Lobby struct {
	ID        uint64
	Engine    engine.Engineer
	Guests    guest.List
	size      int
	isRunning bool
	config    *config.Config
	sync.Mutex
}

// GenerateOneLobby is temporary until lobbies is fully working
// TODO: Remove once lobbies are working
func GenerateOneLobby(config *config.Config) (lobby *Lobby) {
	for _, lobby := range lobbies {
		return lobby
	}

	return NewLobby(config)
}

// NewLobby creates a new lobby with its own engine and list of guests
func NewLobby(config *config.Config) (lobby *Lobby) {
	e := engine.NewEngine(config)
	id := newID()
	lobby = &Lobby{
		ID:     id,
		Engine: e,
		Guests: generalPopulation(lobbySize),
		size:   lobbySize,
		config: config,
	}

	lobbies[lobby.ID] = lobby

	return
}

//This is just for testing until it's fully implemented
const baseID uint64 = 2019968050

func newID() uint64 {
	return baseID
}

func (lobby *Lobby) Start() {
	lobby.Lock()
	defer lobby.Unlock()

	if !lobby.isRunning {
		lobby.isRunning = true
		go lobby.runEngine()
	}
}

// This must be run in a go routine otherwise it will block the thread
func (lobby *Lobby) runEngine() {
	for lobby.Guests.UserSize >= 0 {
		var entities engine.EntityList

		for _, g := range lobby.Guests.Iterator() {
			entity := lobby.Engine.CreateEntity(engine.NewCreature(g, lobby.config))
			entities = append(entities, entity)
		}

		lobby.Engine.Run(entities)
		lobby.Start()
	}

	lobby.isRunning = false
	destroyLobby(lobby)
}

// GetLobby returns the Lobby with id: `id` and if it does not exist it returns
// `(nil, false)`
func GetLobby(id uint64) (*Lobby, bool) {
	l, ok := lobbies[id]
	return l, ok
}

// destroyLobby removes lobby from the global map
func destroyLobby(lobby *Lobby) {
	lobby.Guests.List = nil
	lobby.Guests.UserSize = 0
	//lobby.Engine.Close() - for the runner to be shut down
	lobby.Engine = nil
	delete(lobbies, lobby.ID)
}

func generalPopulation(size int) guest.List {
	var bots guest.List

	for i := size; i > 0; i-- {
		name := fmt.Sprintf("Bot %d", i)
		bots.List = append(bots.List, guest.NewBot(name))
	}

	return bots
}

func (lobby *Lobby) ContainsUser(name string) bool {
	for _, g := range lobby.Guests.Iterator() {
		if g.Name == name {
			return true
		}
	}

	return false
}

// PRE and POST condition for AddUser and RemoveUser:
// The order of the guest list is [0, len-userSize) are botType and
// [len-userSize, len) are userType

// AddUser adds the specified user to the lobby, returning an error if the
// lobby is already at capacity
func (lobby *Lobby) AddUser(user *guest.Guest) error {
	// TODO: Check for duplicates
	lobbyGuests := lobby.Guests.List

	// Check for bot to remove
	var bot *guest.Guest
	bot, lobbyGuests = lobbyGuests[0], lobbyGuests[1:]
	if bot.Type != guest.BotType {
		return errors.New("Lobby full")
	}

	i := len(lobbyGuests) - (lobby.Guests.UserSize + 1)
	newGuest := []*guest.Guest{user}
	lobbyGuests = append(lobbyGuests[:i], append(newGuest, lobbyGuests[i:]...)...)
	lobby.Guests.UserSize++

	lobby.Guests.List = lobbyGuests
	return nil
}

// RemoveUser removes the specified user from the lobby, returning an error if the
// user is not in the lobby
func (lobby *Lobby) RemoveUser(user *guest.Guest) error {
	// TODO: Check for duplicates
	lobbyGuests := lobby.Guests.List

	// Remove User
	var i int
	for i = len(lobbyGuests) - 1; i > 0; i-- {
		if lobbyGuests[i].Name == user.Name {
			// Memory leaks - Go needs to sort slices out...
			copy(lobbyGuests[i:], lobbyGuests[i+1:])
			lobbyGuests[len(lobbyGuests)-1] = nil
			lobbyGuests = lobbyGuests[:len(lobbyGuests)-1]
			break
		}
	}

	// Add Bot
	name := fmt.Sprintf("Bot %d", i+1)
	// This will change in another branch that is getting merged a little later
	newGuest := []*guest.Guest{guest.NewBot(name)}
	lobbyGuests = append(newGuest, lobbyGuests...)
	lobby.Guests.UserSize--

	lobby.Guests.List = lobbyGuests
	return nil
}

// GetLobbyIDs returns an array of all the IDs in the lobbies map
func GetLobbyIDs() []uint64 {
	ids := make([]uint64, 0, len(lobbies))
	for id := range lobbies {
		ids = append(ids, id)
	}
	return ids
}
