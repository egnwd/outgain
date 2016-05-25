package engine

import (
	"github.com/egnwd/outgain/server/protocol"
	"github.com/lucasb-eyer/go-colorful"
	"math"
	"math/rand"
	"time"
)

const gridSize float64 = 10

const defaultRadius float64 = 0.5
const resourceRadius float64 = 0.1

const resourceSpawnInterval time.Duration = 5 * time.Second
const fullUpdateInterval time.Duration = 1 * time.Second

type Engine struct {
	Events <-chan protocol.Event

	eventsOut         chan<- protocol.Event
	tickInterval      time.Duration
	entities          []Entity
	lastTick          time.Time
	lastResourceSpawn time.Time
	nextId            <-chan uint64

	lastFullUpdate time.Time

	worldState protocol.WorldState
}

type Entity interface {
	Tick(dt float64)
	Serialize() protocol.Entity
}

func NewEngine(creatureCount int) (engine *Engine) {
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
		entities:          make([]Entity, 0),
		nextId:            idChannel,
	}

	for i := 0; i < creatureCount; i++ {
		engine.AddEntity(RandomCreature)
	}

	return
}

func (engine *Engine) Run() {
	engine.lastTick = time.Now()
	engine.lastResourceSpawn = time.Now()

	for {
		now := time.Now()

		if now.Sub(engine.lastFullUpdate) > fullUpdateInterval {
			engine.lastFullUpdate = now

			engine.worldState = engine.Serialize()

			engine.eventsOut <- protocol.Event{
				Type: "state",
				Data: engine.worldState,
			}
		} else {
			newState := engine.Serialize()
			diff := protocol.DiffWorld(engine.worldState, newState)
			engine.worldState = newState

			engine.eventsOut <- protocol.Event{
				Type: "diff",
				Data: diff,
			}
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
	engine.entities = append(engine.entities, entity)
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
}

type Creature struct {
	Id     uint64
	Name   string
	Color  string
	X      float64
	Y      float64
	Radius float64

	dx float64
	dy float64
}

func RandomCreature(id uint64) Entity {
	angle := rand.Float64() * 2 * math.Pi
	x := rand.Float64() * gridSize
	y := rand.Float64() * gridSize

	return &Creature{
		Id:   id,
		Name: "foo",

		Color:  colorful.FastHappyColor().Hex(),
		X:      x,
		Y:      y,
		Radius: defaultRadius,
		dx:     math.Cos(angle),
		dy:     math.Sin(angle),
	}
}

func (creature *Creature) Tick(dt float64) {
	angle := rand.NormFloat64() * math.Pi / 4
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	dx := creature.dx*cos - creature.dy*sin
	dy := creature.dx*sin + creature.dy*cos
	creature.dx = dx
	creature.dy = dy

	creature.X += creature.dx * dt
	creature.Y += creature.dy * dt

	if creature.X-creature.Radius < 0 {
		creature.X = creature.Radius
		creature.dx *= -1
	}
	if creature.X+creature.Radius > gridSize {
		creature.X = gridSize - creature.Radius
		creature.dx *= -1
	}
	if creature.Y-creature.Radius < 0 {
		creature.Y = creature.Radius
		creature.dy *= -1
	}
	if creature.Y+creature.Radius > gridSize {
		creature.Y = gridSize - creature.Radius
		creature.dy *= -1
	}
}

func (creature *Creature) Serialize() protocol.Entity {
	return protocol.Entity{
		Id:     creature.Id,
		Name:   &creature.Name,
		Color:  creature.Color,
		X:      creature.X,
		Y:      creature.Y,
		Radius: creature.Radius,
	}
}

type Resource struct {
	Id     uint64
	Color  string
	X      float64
	Y      float64
	Radius float64
}

func RandomResource(id uint64) Entity {
	return &Resource{
		Id:     id,
		X:      rand.Float64() * gridSize,
		Y:      rand.Float64() * gridSize,
		Radius: resourceRadius,
		Color:  colorful.FastHappyColor().Hex(),
	}
}

func (resource *Resource) Tick(dt float64) {
}

func (resource *Resource) Serialize() protocol.Entity {
	return protocol.Entity{
		Id:     resource.Id,
		Name:   nil,
		Color:  resource.Color,
		X:      resource.X,
		Y:      resource.Y,
		Radius: resource.Radius,
	}
}
