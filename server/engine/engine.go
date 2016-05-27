package engine

import (
	"time"

	"github.com/egnwd/outgain/server/protocol"
)

const gridSize float64 = 10

const resourceSpawnInterval time.Duration = 5 * time.Second

const eatRadiusDifference = 0.2

const initialCreatureCount = 10

type Engine struct {
	Events <-chan protocol.Event

	eventsOut         chan<- protocol.Event
	tickInterval      time.Duration
	entities          EntityList
	lastTick          time.Time
	lastResourceSpawn time.Time
	nextId            <-chan uint64
}

func NewEngine() (engine *Engine) {
	eventChannel := make(chan protocol.Event)
	idChannel := make(chan uint64)
	go func() {
		var id uint64 = 0
		for {
			idChannel <- id
			id += 1
		}
	}()

	engine = &Engine{
		Events:            eventChannel,
		eventsOut:         eventChannel,
		tickInterval:      time.Millisecond * 100,
		lastTick:          time.Now(),
		lastResourceSpawn: time.Now(),
		entities:          EntityList{},
		nextId:            idChannel,
	}

	return
}

func (engine *Engine) Reset() {
	engine.entities = EntityList{}

	for i := 0; i < initialCreatureCount; i++ {
		engine.AddEntity(RandomCreature)
	}
	engine.clearGameLog()
}

// clearGameLog should clear the current game-log (or make it clear that a new game has begun)
func (engine *Engine) clearGameLog() {
	logEvent := protocol.LogEvent{0, 0, 0}
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

func (engine *Engine) AddEntity(builder func(uint64) Entity) {
	entity := builder(<-engine.nextId)
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
		logEvent = protocol.LogEvent{1, a.Base().Id, 0}
	case *Creature:
		logEvent = protocol.LogEvent{2, a.Base().Id, b.Base().Id}
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

func (engine *Engine) collisionDetection() {
	for _, entity := range engine.entities {
		entity.Base().dying = false
		entity.Base().radiusIncrement = 0
	}

	for collision := range engine.entities.Collisions() {
		a, b := collision.a, collision.b
		diff := a.Base().Radius - b.Base().Radius
		if diff > eatRadiusDifference {
			a.Base().radiusIncrement += b.Base().Radius + b.Base().radiusIncrement
			b.Base().dying = true
			engine.addLogEvent(a, b)
		} else if diff < -eatRadiusDifference {
			b.Base().radiusIncrement += a.Base().Radius + a.Base().radiusIncrement
			a.Base().dying = true
			engine.addLogEvent(b, a)
		}
	}

	engine.entities = engine.entities.Filter(func(entity Entity) bool {
		return !entity.Base().dying
	})

	var resetEngine = false
	for _, entity := range engine.entities {
		entity.Base().Radius += entity.Base().radiusIncrement

		if entity.Base().Radius > gridSize/2 {
			resetEngine = true
		}
	}

	// Changing the radius of entities changes their left coordinate,
	// so sort the list again to maintain the invariant
	engine.entities.Sort()

	if resetEngine {
		engine.Reset()
	}
}
