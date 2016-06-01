package lobby

import (
	"errors"
	"fmt"
	"log"

	"github.com/egnwd/outgain/server/engine"
	"github.com/egnwd/outgain/server/guest"
)

const lobbySize int = 10

var lobbies = make(map[uint64]*Lobby)

type guestList struct {
	list     []*guest.Guest
	userSize int
}

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
		lobby.Engine.AddEntity(guest, engine.RandomCreature)
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
	lobby.Guests.list = []*guest.Guest{}
	lobby.Guests.userSize = 0
	lobby.Engine = nil
	delete(lobbies, lobby.ID)
}

func generalPopulation(size int) guestList {
	var bots guestList

	for i := size; i > 0; i-- {
		name := fmt.Sprintf("Bot %d", i)
		bots.list = append(bots.list, guest.NewBot(name))
	}

	return bots
}

func (lobby *Lobby) ContainsUser(name string) bool {
	for _, g := range lobby.Guests.list {
		if g.Name == name {
			return true
		}
	}

	return false
}

// AddUser adds the specified user to the lobby, returning an error if the
// lobby is already at capacity, and running the engine if the user is
// the first to join
func (lobby *Lobby) AddUser(user *guest.Guest) error {
	// TODO: Assert User
	lobbyGuests := lobby.Guests.list

	// Check for bot to remove
	bot, newGuests := lobbyGuests[0], lobbyGuests[1:]
	if bot.Type != guest.BotType {
		return errors.New("Lobby full")
	}

	i := len(lobbyGuests) - (lobby.Guests.userSize + 1)
	newGuest := []*guest.Guest{user}
	newGuests = append(newGuests[:i], append(newGuest, newGuests[i:]...)...)
	lobby.Guests.userSize++

	log.Printf("%d\n", lobby.Guests.userSize)

	lobby.Guests.list = newGuests
	if lobby.Guests.userSize == 1 {
		lobby.startEngine()
	}
	return nil
}
