package lobby

import (
	"encoding/json"
	"log"
	"math/rand"
	"sync"

	"github.com/egnwd/outgain/server/config"
	"github.com/egnwd/outgain/server/engine"
	"gopkg.in/antage/eventsource.v1"
)

const lobbySize int = 10

var lobbies = make(map[uint64]*Lobby)

// Lobby runs its own instance of an engine, and keeps track of its users
type Lobby struct {
	ID        uint64
	Engine    engine.Engineer
	Events    eventsource.EventSource
	Guests    guestList
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
		Guests: generalPopulation(lobbySize),
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

	for lobby.Guests.userSize > 0 {
		var entities engine.EntityList

		for _, guest := range lobby.Guests.list {
			entity := lobby.Engine.CreateEntity(engine.NewCreature(guest.name, lobby.config))
			entities = append(entities, entity)
		}

		lobby.Engine.Run(entities)
		log.Println("Finished Running")
		log.Printf("Users in Game: %d\n", lobby.Guests.userSize)
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
	lobby.Guests.list = nil
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
