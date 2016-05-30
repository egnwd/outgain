package engine

import (
	"math"
	"time"

	"github.com/egnwd/outgain/server/config"
	"github.com/egnwd/outgain/server/protocol"
)

const gridSize float64 = 10

const resourceSpawnInterval time.Duration = 5 * time.Second

const eatRadiusDifference = 0.2

const initialCreatureCount = 10

// Engine stores the information about an instance of the game and controls
// the events that are occuring within the game
type Engine struct {
	Events <-chan protocol.Event

	eventsOut         chan<- protocol.Event
	tickInterval      time.Duration
	entities          EntityList
	lastTick          time.Time
	lastResourceSpawn time.Time
	nextID            <-chan uint64
	config            *config.Config
}

// NewEngine returns a fresh instance of a game engine
func NewEngine(config *config.Config) (engine *Engine) {
	eventChannel := make(chan protocol.Event)
	idChannel := make(chan uint64)
	go func() {
		var id uint64 = 0
		for {
			idChannel <- id
			id++
		}
	}()

	engine = &Engine{
		Events:            eventChannel,
		eventsOut:         eventChannel,
		tickInterval:      time.Millisecond * 100,
		lastTick:          time.Now(),
		lastResourceSpawn: time.Now(),
		entities:          EntityList{},
		nextID:            idChannel,
		config:            config,
	}

	return
}

// Reset clears the game log and resets the creatures
func (engine *Engine) Reset() {
	for _, entity := range engine.entities {
		entity.Close()
	}

	engine.entities = EntityList{}

	for i := 0; i < initialCreatureCount; i++ {
		engine.AddEntity(RandomCreature(engine.config))
	}
	engine.clearGameLog()
}

// clearGameLog should clear the current game-log (or make it clear that a new game has begun)
func (engine *Engine) clearGameLog() {
	logEvent := protocol.LogEvent{LogType: 0, ProtagID: 0, AntagID: 0}
	engine.eventsOut <- protocol.Event{
		Type: "log",
		Data: logEvent,
	}
}

func (engine *Engine) Run() {
	engine.Reset()
	engine.lastTick = time.Now()
	engine.lastResourceSpawn = time.Now()

	for {
		engine.eventsOut <- protocol.Event{
			Type: "state",
			Data: engine.Serialize(),
		}

		time.Sleep(engine.tickInterval)

		engine.tick()
	}
}

func (engine *Engine) Serialize() protocol.WorldState {
	entities := make([]protocol.Entity, len(engine.entities))
	for i, entity := range engine.entities {
		entities[i] = entity.Serialize()
	}

	return protocol.WorldState{
		Time:     uint64(engine.lastTick.UnixNano()) / 1e6,
		Entities: entities,
	}
}

// AddEntity adds an entity to the engine's list
func (engine *Engine) AddEntity(builder func(uint64) Entity) {
	entity := builder(<-engine.nextID)
	engine.entities = engine.entities.Insert(entity)
}

// addLogEvent adds to  logEvents which are eventually added to the gameLog
// Where is the best place to document the number -> eventType mappings?
func (engine *Engine) addLogEvent(a, b Entity) {
	var logEvent protocol.LogEvent
	switch b.(type) {
	case nil:
		return
	case *Resource:
		logEvent = protocol.LogEvent{LogType: 1, ProtagID: a.Base().ID, AntagID: 0}
	case *Creature:
		logEvent = protocol.LogEvent{LogType: 2, ProtagID: a.Base().ID, AntagID: b.Base().ID}
	}
	engine.eventsOut <- protocol.Event{
		Type: "log",
		Data: logEvent,
	}

}

func (engine *Engine) tick() {
	now := time.Now()
	dt := now.Sub(engine.lastTick).Seconds()
	engine.lastTick = now

	if now.Sub(engine.lastResourceSpawn) > resourceSpawnInterval {
		engine.lastResourceSpawn = now

		engine.AddEntity(RandomResource)
	}

	engine.entities.Tick(dt)
	engine.collisionDetection()
}

func (engine *Engine) eatEntity(eater, eaten Entity) {
	eater.Base().nextRadius = math.Sqrt(eater.Volume() + eaten.Volume())
	eaten.Base().dying = true
	engine.addLogEvent(eater, eaten)
}

func (engine *Engine) collisionDetection() {
	for _, entity := range engine.entities {
		entity.Base().dying = false
		entity.Base().nextRadius = entity.Base().Radius
	}

	for collision := range engine.entities.Collisions() {
		a, b := collision.a, collision.b
		diff := a.Base().Radius - b.Base().Radius

		if diff >= eatRadiusDifference {
			engine.eatEntity(a, b)
		} else if diff <= -eatRadiusDifference {
			engine.eatEntity(b, a)
		}
	}

	engine.entities = engine.entities.Filter(func(entity Entity) bool {
		return !entity.Base().dying
	})

	creatureCount := 0
	for _, entity := range engine.entities {
		entity.Base().Radius = entity.Base().nextRadius

		_, isCreature := entity.(*Creature)
		if isCreature {
			creatureCount++
		}
	}

	// Changing the radius of entities changes their left coordinate,
	// so sort the list again to maintain the invariant
	engine.entities.Sort()

	if creatureCount <= 1 {
		engine.Reset()
	}
}
