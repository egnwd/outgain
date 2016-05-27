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
	logEvents         []LogEvent
	tickInterval      time.Duration
	entities          EntityList
	lastTick          time.Time
	lastResourceSpawn time.Time
	nextId            <-chan uint64
}

type LogEvent struct {
	logType  int    // e.g. 0 for new game, 1 for creature -> resource etc
	protagID uint64 // ID of the protagonist (if any)
	antagID  uint64 // ID of the antagonist (if any)
}

// Refactor?
func (logEvent *LogEvent) Serialize() protocol.LogEvent {
	return protocol.LogEvent{
		LogType:  logEvent.logType,
		ProtagID: logEvent.protagID,
		AntagID:  logEvent.antagID,
	}
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
	clearGameLog(&engine.logEvents)
}

// clearGameLog should clear the current game-log (or make it clear that a new game has begun)
func clearGameLog(logEvents *[]LogEvent) {
	*logEvents = append(*logEvents, LogEvent{0, 0, 0})
}

func (engine *Engine) Run() {
	engine.lastTick = time.Now()
	engine.lastResourceSpawn = time.Now()

	for {
		engine.eventsOut <- protocol.Event{
			Type: "state",
			Data: engine.Serialize(),
		}

		engine.eventsOut <- protocol.Event{
			Type: "log",
			Data: engine.SerializeLog(),
		}

		engine.logEvents = []LogEvent{}

		time.Sleep(engine.tickInterval)

		engine.tick()
	}
}

func (engine *Engine) SerializeLog() protocol.LogEvents {
	lEvents := make([]protocol.LogEvent, len(engine.logEvents))
	for i, logEvent := range engine.logEvents {
		lEvents[i] = logEvent.Serialize()
	}
	return protocol.LogEvents{
		LogEvents: lEvents,
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
func addLogEvent(logEvents *[]LogEvent, a, b Entity) {
	switch b.(type) {
	case nil:
		// I don't know Go well enough to know what to put here, open to suggestions
	case *Resource:
		*logEvents = append(*logEvents, LogEvent{1, a.Base().Id, 0})
	case *Creature:
		*logEvents = append(*logEvents, LogEvent{2, a.Base().Id, b.Base().Id})
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
			addLogEvent(&engine.logEvents, a, b)
		} else if diff < -eatRadiusDifference {
			b.Base().radiusIncrement += a.Base().Radius + a.Base().radiusIncrement
			a.Base().dying = true
			addLogEvent(&engine.logEvents, b, a)
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
