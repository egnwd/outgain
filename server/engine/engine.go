package engine

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/egnwd/outgain/server/database"
	"github.com/egnwd/outgain/server/protocol"
)

const gridSize float64 = 10

const baseResourceSpawnInterval = 5
const maxResourceIncrease = 5

var resourceSpawnInterval time.Duration

const eatRadiusDifference = 0.1

const initialCreatureCount = 10
const drainRate = 0.5
const radiusThreshold = 0.2

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
	restart           bool
}

type builderFunc func(uint64) Entity

// NewEngine returns a fresh instance of a game engine
func NewEngine() (engine *Engine) {
	eventChannel := make(chan protocol.Event)
	idChannel := make(chan uint64)
	go func() {
		var id uint64
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
	}

	return
}

// restartEngine puts the engine back to it's original state
func (engine *Engine) restartEngine() {
	engine.updateLeaderboard()
	for _, entity := range engine.entities {
		entity.Close()
	}

	engine.entities = EntityList{}
	engine.clearGameLog()

	log.Println("Restarting Engine")
	engine.restart = true
}

func (engine *Engine) updateLeaderboard() {
	//users := engine.entities.Filter(func(entity Entity) bool {
	//	creature, isCreature := entity.(*Creature)
	//	if isCreature {
	//		return creature.Guest.Type == guest.UserType
	//	}
	//	return false
	//})
	//users.SortScore()
	//for _, user := range users {
	//	fmt.Println(user.GetGains())
	//}
	engine.entities.SortScore()
	for _, entity := range engine.entities {
		var minVal = database.GetMinScore()
		if gains := entity.GetGains(); gains > minVal {
			database.UpdateLeaderboard(entity.GetName(), gains)
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
	engine.clearGameLog()
	engine.lastTick = time.Now()
	regenerateResourceInterval()
	// engine.lastResourceSpawn = time.Now()

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

		if engine.restart {
			engine.restart = false
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
	var logEvent protocol.LogEvent
	switch b.(type) {
	case nil:
		return
	case *Resource:
		logEvent = protocol.LogEvent{
			LogType:    1,
			ProtagName: a.GetName(),
			AntagName:  b.GetName(),
			Gains:      a.GetGains(),
		}
	case *Creature:
		logEvent = protocol.LogEvent{
			LogType:    2,
			ProtagName: a.GetName(),
			AntagName:  b.GetName(),
			Gains:      a.GetGains(),
		}
	case *Spike:
		logEvent = protocol.LogEvent{
			LogType:    3,
			ProtagName: a.GetName(),
			AntagName:  b.GetName(),
			Gains:      0,
		}
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
	engine.collisionDetection(dt)
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
	if eatenIsSpike {
		if eater.Base().nextRadius <= defaultRadius {
			eater.Base().dying = true
		}
		eater.Base().nextRadius = math.Sqrt(eaterVolume / 2)
		eaten.Base().dying = true
	} else {

		amount := math.Exp(-1/drainRate*dt) * eatenVolume

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

	// We currently run both the slow and fast collision algorithms, and
	// compare their outputs to find collision missed by the fast one due
	// to bugs. Once we're confident enough with the results of the fast
	// one we can switch fully to this one.
	collisions := []Collision{}
	for collision := range engine.entities.Collisions() {
		collisions = append(collisions, collision)

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

	for collision := range engine.entities.SlowCollisions() {
		found := false
		for _, c := range collisions {
			if (c.A == collision.A && c.B == collision.B) ||
				(c.A == collision.B && c.B == collision.A) {
				found = true
				break
			}
		}

		if !found {
			message := fmt.Sprintf("WARN Collision false negative: (%d %d),",
				collision.A.Base().ID, collision.B.Base().ID)

			for _, e := range engine.entities {
				message += fmt.Sprintf(" dummyEntity(%d, %.2f, %.2f, %.2f),",
					e.Base().ID, e.Base().X, e.Base().Y, e.Base().Radius)
			}

			log.Println(message)
		}
	}

	engine.entities = engine.entities.Filter(func(entity Entity) bool {
		return !entity.Base().dying
	})

	for _, entity := range engine.entities {
		if entity.Base().dying {
			fmt.Println("Paul is a bad programmer")
		}
	}

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
		engine.restartEngine()
	}
}
