package engine

import (
	"fmt"
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
	events            []string
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

	engine.Reset()
	return
}

func (engine *Engine) Reset() {
	engine.entities = EntityList{}

	for i := 0; i < initialCreatureCount; i++ {
		engine.AddEntity(RandomCreature)
	}
	clearGameLog(&engine.events)
}

// clearGameLog should clear the current game-log (or make it clear that a new game has begun
func clearGameLog(events *[]string) {
	// This is temporary before I deal with the front-end
	message := fmt.Sprintf("A new game has begun!\n\n\n")
	*events = append(*events, message)
}

func (engine *Engine) Run() {
	engine.lastTick = time.Now()
	engine.lastResourceSpawn = time.Now()

	for {
		engine.eventsOut <- protocol.Event{
			Type: "state",
			Data: engine.Serialize(),
		}

		engine.events = []string{}

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
		Time:      uint64(engine.lastTick.UnixNano()) / 1e6,
		Entities:  entities,
		LogEvents: engine.events,
	}
}

func (engine *Engine) AddEntity(builder func(uint64) Entity) {
	entity := builder(<-engine.nextId)
	engine.entities = engine.entities.Insert(entity)
}

// addEvent adds events which are eventually added to the gameLog
// TODO: Refactor to send minimal information with inference on client-side
func addEvent(events *[]string, a, b Entity) {
	switch b.(type) {
	case nil:
		// I don't know Go well enough to know what to put here, open to suggestions
	case *Resource:
		message := fmt.Sprintf("Yum, creature %d ate a resource\n", a.Base().Id)
		*events = append(*events, message)
	case *Creature:
		message := fmt.Sprintf("Creature number %d ate creature %d\n", a.Base().Id, b.Base().Id)
		*events = append(*events, message)
	}

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

	message := fmt.Sprintf("Test - %s\n", now.String())
	engine.events = append(engine.events, message)
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
			addEvent(&engine.events, a, b)
		} else if diff < -eatRadiusDifference {
			b.Base().radiusIncrement += a.Base().Radius + a.Base().radiusIncrement
			a.Base().dying = true
			addEvent(&engine.events, b, a)
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
