package engine

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/egnwd/outgain/server/database"
	"github.com/egnwd/outgain/server/protocol"
)

const gridSize float64 = 20

const baseResourceSpawnInterval = 5
const maxResourceIncrease = 5

var resourceSpawnInterval time.Duration

const eatRadiusDifference = 0.1

const initialCreatureCount = 10
const drainRate = 0.5
const radiusThreshold = 0.2

const roundLength = 15 * time.Second

// Engine stores the information about an instance of the game and controls
// the events that are occuring within the game
type Engine struct {
	eventsOut         chan<- protocol.Event
	tickInterval      time.Duration
	entities          EntityList
	users             EntityList
	lastTick          time.Time
	lastResourceSpawn time.Time
	nextID            <-chan uint64
	restarted         bool
}

type builderFunc func(uint64) Entity

// NewEngine returns a fresh instance of a game engine
func NewEngine(eventChannel chan protocol.Event) (engine *Engine) {
	idChannel := make(chan uint64)
	go func() {
		var id uint64
		for {
			idChannel <- id
			id++
		}
	}()

	engine = &Engine{
		eventsOut:         eventChannel,
		tickInterval:      time.Millisecond * 100,
		lastTick:          time.Now(),
		lastResourceSpawn: time.Now(),
		entities:          EntityList{},
		users:             EntityList{},
		nextID:            idChannel,
	}

	return
}

// restart puts the engine back to it's original state
func (engine *Engine) restart() {
	engine.updateLeaderboard()
	for _, entity := range engine.entities {
		entity.Close()
	}

	engine.entities = EntityList{}
	engine.users = EntityList{}

	engine.eventsOut <- protocol.Event{
		Type: "state",
		Data: engine.Serialize(),
	}

	engine.clearGameLog()

	log.Println("Restarting Engine")
	engine.restarted = true
}

func (engine *Engine) updateLeaderboard() {
	engine.users = engine.users.SortScore()

	for _, entity := range engine.users {
		var minVal = database.GetMinScore()
		if gains := entity.GetGains(); gains > minVal {
			if entity.IsUser() { // Bots can't set high scores
				database.UpdateLeaderboard(entity.GetName(), gains)
			}
		} else {
			break // The list is sorted, no need to check the rest
		}
	}

}

// clearGameLog should clear the current game-log (or make it clear that a new game has begun)
func (engine *Engine) clearGameLog() {
	logEvent := protocol.LogEvent{LogType: 0, ProtagName: "", AntagName: ""}
	engine.eventsOut <- protocol.Event{
		Type: "log",
		Data: logEvent,
	}
}

// Run starts the simulation of the game
func (engine *Engine) Run(entities EntityList) {
	log.Println("Running Engine")
	engine.entities = entities
	for _, entity := range entities {
		if entity.IsUser() {
			engine.users = append(engine.users, entity)
		}
	}
	engine.clearGameLog()
	engine.lastTick = time.Now()
	regenerateResourceInterval()

	roundTimer := time.NewTimer(roundLength)

GameLoop:
	for {
		engine.eventsOut <- protocol.Event{
			Type: "state",
			Data: engine.Serialize(),
		}

		wakeup := engine.lastTick.Add(engine.tickInterval)
		if now := time.Now(); wakeup.After(now) {
			time.Sleep(wakeup.Sub(now))
		}

		engine.tick()

		select {
		case <-roundTimer.C:
			engine.restart()
		default:
		}

		if engine.restarted {
			engine.restarted = false
			break GameLoop
		}

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
func (engine *Engine) AddEntity(builder builderFunc) {
	engine.entities = engine.entities.Insert(engine.CreateEntity(builder))
}

// CreateEntity builds an entity using the builder
func (engine *Engine) CreateEntity(builder builderFunc) Entity {
	entity := builder(<-engine.nextID)
	switch entity.(type) {
	case *Spike:
		for {
			collides := false
			for _, candidateCreature := range engine.entities {
				if Collide(entity, candidateCreature) {
					collides = true
					break
				}
			}

			if !collides {
				break
			}
			entity = builder(<-engine.nextID)
		}
	}
	return entity
}

// addLogEvent adds to  logEvents which are eventually added to the gameLog
// Where is the best place to document the number -> eventType mappings?
func (engine *Engine) addLogEvent(a, b Entity) {
	var (
		logEvent protocol.LogEvent
		logType  int
	)
	switch b.(type) {
	case nil:
		return
	case *Resource:
		logType = 1
		break
	case *Creature:
		logType = 2
		break
	case *Spike:
		logType = 3
	}

	logEvent = protocol.LogEvent{
		LogType:    logType,
		ProtagName: a.GetName(),
		AntagName:  b.GetName(),
		Gains:      a.GetGains(),
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

		regenerateResourceInterval()
		for i := -1; i < rand.Intn(maxResourceIncrease); i++ {
			engine.AddEntity(RandomResource)
		}
		engine.AddEntity(RandomSpike)
	}

	state := engine.Serialize()

	engine.entities.Tick(state, dt)
	engine.entities = engine.entities.Filter(func(entity Entity) bool {
		return !entity.Base().dying
	})

	engine.collisionDetection(dt)
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

	if creatureCount <= 1 {
		engine.restartEngine()
	}
}

func regenerateResourceInterval() {
	rand.Seed(time.Now().Unix())
	duration := time.Duration(rand.Intn(baseResourceSpawnInterval) + 1)
	resourceSpawnInterval = duration * time.Second
}

func (engine *Engine) eatEntity(dt float64, eater, eaten Entity) {
	_, eaterIsResource := eater.(*Resource)
	if eaterIsResource || eater.Base().dying || eaten.Base().dying {
		return
	}
	_, eaterIsSpike := eater.(*Spike)
	_, eatenIsSpike := eaten.(*Spike)
	if eaterIsSpike {
		engine.eatEntity(dt, eaten, eater)
		return
	}

	eaterVolume := eater.Base().nextRadius * eater.Base().nextRadius
	eatenVolume := eaten.Base().nextRadius * eaten.Base().nextRadius

	var amount float64

	if eatenIsSpike {
		amount = math.Exp(-1/drainRate*dt) * eaterVolume
		nextCreatureVolume := eaterVolume + amount*eaten.BonusFactor()
		if nextCreatureVolume < 0 {
			eater.Base().dying = true
		} else {
			eater.Base().nextRadius = math.Sqrt(nextCreatureVolume)
		}
		eaten.Base().dying = true
		eater.(*Creature).decrementScore()
	} else {
		amount = math.Exp(-1/drainRate*dt) * eatenVolume
		eater.Base().nextRadius = math.Sqrt(eaterVolume + amount*eaten.BonusFactor())
		eaten.Base().nextRadius = math.Sqrt(eatenVolume - amount)
		if eaten.Base().nextRadius < radiusThreshold {
			eater.(*Creature).incrementScore(eaten)
			eaten.Base().dying = true
		}
	}

	engine.addLogEvent(eater, eaten)
}

func (engine *Engine) collisionDetection(dt float64) {
	for _, entity := range engine.entities {
		entity.Base().dying = false
		entity.Base().nextRadius = entity.Base().Radius
	}

	for collision := range engine.entities.Collisions() {
		a, b := collision.A, collision.B
		diff := a.Base().Radius - b.Base().Radius

		if diff >= eatRadiusDifference {
			engine.eatEntity(dt, a, b)
		} else if diff <= -eatRadiusDifference {
			engine.eatEntity(dt, b, a)
		} else {
			switch a.(type) {
			case *Spike:
				engine.eatEntity(dt, b, a)
			}

			switch b.(type) {
			case *Spike:
				engine.eatEntity(dt, a, b)
			}
		}
	}

	for _, entity := range engine.entities {
		entity.Base().Radius = entity.Base().nextRadius
	}

	// Changing the radius of entities changes their left coordinate,
	// so sort the list again to maintain the invariant
	engine.entities.SortLeft()
}
