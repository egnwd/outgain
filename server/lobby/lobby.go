package lobby

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"

	"sync"

	"github.com/egnwd/outgain/server/config"
	"github.com/egnwd/outgain/server/engine"
	"github.com/egnwd/outgain/server/guest"
	"gopkg.in/antage/eventsource.v1"
)

const lobbySize int = 10

var lobbies = make(map[uint64]*Lobby)

// Lobby runs its own instance of an engine, and keeps track of its users
type Lobby struct {
	ID        uint64
	Engine    engine.Engineer
	Events    eventsource.EventSource
	Guests    guest.List
	size      int
	isRunning bool
	config    *config.Config
	sync.Mutex
}

// NewLobby creates a new lobby with its own engine and list of guests
func NewLobby(config *config.Config) (lobby *Lobby) {
	engine := engine.NewEngine()
	events := eventsource.New(nil, nil)
	id := newID()
	lobby = &Lobby{
		ID:     id,
		Engine: engine,
		Events: events,
		Guests: generalPopulation(lobbySize, config),
		size:   lobbySize,
		config: config,
	}

	go func() {
		for event := range engine.Events {
			packet, err := json.Marshal(event.Data)
			if err != nil {
				log.Printf("JSON serialization failed %v", err)
			} else {
				events.SendEventMessage(string(packet), event.Type, "")
			}
		}
	}()

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
	log.Println("Running game in lobby")

	for lobby.Guests.UserSize > 0 {
		var entities engine.EntityList

		for _, g := range lobby.Guests.Iterator() {
			entity := lobby.Engine.CreateEntity(engine.NewCreature(g, lobby.config))
			entities = append(entities, entity)
		}

		lobby.Engine.Run(entities)
		log.Println("Finished Running")
		log.Printf("Users in Game: %d\n", lobby.Guests.UserSize)
	}

	log.Println("Destroying Lobby")
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

func generalPopulation(size int, config *config.Config) guest.List {
	var bots guest.List

	source, err := ioutil.ReadFile(config.DefaultAI)
	if err != nil {
		log.Fatalln(err)
	}

	for i := size; i > 0; i-- {
		name := fmt.Sprintf("Bot %d", i)
		bots.List = append(bots.List, guest.NewBot(name, string(source)))
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
func (lobby *Lobby) AddUser(username string) error {
	// TODO: Check for duplicates
	lobbyGuests := lobby.Guests.List

	// Check for bot to remove
	var bot *guest.Guest
	bot, lobbyGuests = lobbyGuests[0], lobbyGuests[1:]
	if bot.Type != guest.BotType {
		return errors.New("Lobby full")
	}

	source, err := ioutil.ReadFile(lobby.config.DefaultAI)
	if err != nil {
		log.Fatalln(err)
	}
	user := guest.NewUser(username, string(source))

	i := len(lobbyGuests) - lobby.Guests.UserSize
	newGuest := []*guest.Guest{user}
	lobbyGuests = append(lobbyGuests[:i], append(newGuest, lobbyGuests[i:]...)...)
	lobby.Guests.UserSize++

	lobby.Guests.List = lobbyGuests
	return nil
}

// RemoveUser removes the specified user from the lobby, returning an error if the
// user is not in the lobby
func (lobby *Lobby) RemoveUser(username string) error {
	// TODO: Check for duplicates
	lobbyGuests := lobby.Guests.List

	// Remove User
	var i int
	for i = len(lobbyGuests) - 1; i > 0; i-- {
		if lobbyGuests[i].Name == username {
			// Memory leaks - Go needs to sort slices out...
			copy(lobbyGuests[i:], lobbyGuests[i+1:])
			lobbyGuests[len(lobbyGuests)-1] = nil
			lobbyGuests = lobbyGuests[:len(lobbyGuests)-1]
			break
		}
	}

	// Add Bot
	name := fmt.Sprintf("Bot %d", i+1)
	source, err := ioutil.ReadFile(lobby.config.DefaultAI)
	if err != nil {
		log.Fatalln(err)
	}

	// This will change in another branch that is getting merged a little later
	newGuest := []*guest.Guest{guest.NewBot(name, string(source))}
	lobbyGuests = append(newGuest, lobbyGuests...)
	lobby.Guests.UserSize--

	lobby.Guests.List = lobbyGuests
	return nil
}

func (lobby *Lobby) FindGuest(username string) *guest.Guest {
	for _, user := range lobby.Guests.Iterator() {
		if user.Name == username {
			return user
		}
	}

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
