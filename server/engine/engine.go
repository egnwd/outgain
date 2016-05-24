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

type Engine struct {
	Updates <-chan protocol.WorldUpdate

	updatesOut   chan<- protocol.WorldUpdate
	tickInterval time.Duration
	creatures    []*Creature
	lastTick     uint64
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

func NewEngine(creatureCount int) *Engine {
	colors := colorful.FastHappyPalette(creatureCount)

	creatures := make([]*Creature, creatureCount)
	for i := range creatures {
		angle := rand.Float64() * 2 * math.Pi
		x := rand.Float64() * gridSize
		y := rand.Float64() * gridSize

		creatures[i] = &Creature{
			Id:   uint64(i),
			Name: "foo",

			Color:  colors[i].Hex(),
			X:      x,
			Y:      y,
			Radius: defaultRadius,
			dx:     math.Cos(angle),
			dy:     math.Sin(angle),
		}
	}

	ch := make(chan protocol.WorldUpdate)

	return &Engine{
		Updates:      ch,
		updatesOut:   ch,
		tickInterval: time.Millisecond * 100,
		creatures:    creatures,
		lastTick:     0,
	}
}

func (engine *Engine) Run() {
	engine.lastTick = uint64(time.Now().UnixNano()) / 1e6

	for {
		update := protocol.WorldUpdate{
			Time:      engine.lastTick,
			Creatures: make([]protocol.Creature, len(engine.creatures)),
		}

		for i, c := range engine.creatures {
			update.Creatures[i] = protocol.Creature{
				Id:     c.Id,
				Name:   c.Name,
				Color:  c.Color,
				X:      c.X,
				Y:      c.Y,
				Radius: c.Radius,
			}
		}

		engine.updatesOut <- update

		time.Sleep(engine.tickInterval)

		engine.tick()
	}
}

func (engine *Engine) tick() {
	now := uint64(time.Now().UnixNano()) / 1e6
	dt := float64(now-engine.lastTick) / 1000
	engine.lastTick = now

	for _, c := range engine.creatures {
		angle := rand.NormFloat64() * math.Pi / 4
		cos := math.Cos(angle)
		sin := math.Sin(angle)

		dx := c.dx*cos - c.dy*sin
		dy := c.dx*sin + c.dy*cos
		c.dx = dx
		c.dy = dy

		c.X += c.dx * dt
		c.Y += c.dy * dt

		if c.X-c.Radius < 0 {
			c.X = c.Radius
			c.dx *= -1
		}
		if c.X+c.Radius > gridSize {
			c.X = gridSize - c.Radius
			c.dx *= -1
		}
		if c.Y-c.Radius < 0 {
			c.Y = c.Radius
			c.dy *= -1
		}
		if c.Y+c.Radius > gridSize {
			c.Y = gridSize - c.Radius
			c.dy *= -1
		}
	}
}
