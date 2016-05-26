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

func (engine *Engine) tick() {
	now := time.Now()
	dt := now.Sub(engine.lastTick).Seconds()
	engine.lastTick = now

	if now.Sub(engine.lastResourceSpawn) > resourceSpawnInterval {
		engine.lastResourceSpawn = now

		engine.AddEntity(RandomResource)
	}

	for _, entity := range engine.entities {
		entity.Tick(dt)
	}

	engine.entities.Sort()

	for collision := range engine.entities.Collisions() {
		a, b := collision.a, collision.b
		diff := a.Base().Radius - b.Base().Radius
		if diff > eatRadiusDifference {
			a.Base().Radius += b.Base().Radius
			b.Base().Dying = true
		} else if diff < -eatRadiusDifference {
			b.Base().Radius += a.Base().Radius
			a.Base().Dying = true
		}
	}

	engine.entities = engine.entities.Filter(func(entity Entity) bool {
		return !entity.Base().Dying
	})

	for _, entity := range engine.entities {
		if entity.Base().Radius > gridSize/2 {
			engine.Reset()
			break
		}
	}

	message := fmt.Sprintf("Test - %s\n", now.String())
	engine.events = append(engine.events, message)
}
