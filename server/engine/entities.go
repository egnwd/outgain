package engine

import (
	"github.com/egnwd/outgain/server/protocol"
	"github.com/lucasb-eyer/go-colorful"
	"math"
	"math/rand"
)

const defaultRadius float64 = 0.5
const resourceRadius float64 = 0.1

type Entity interface {
	Tick(dt float64)
	Serialize() protocol.Entity
	Base() *EntityBase
}

type EntityBase struct {
	Id     uint64
	Color  string
	X      float64
	Y      float64
	Radius float64
}

type Creature struct {
	EntityBase

	Name string

	dx float64
	dy float64
}

func RandomCreature(id uint64) Entity {
	angle := rand.Float64() * 2 * math.Pi
	x := rand.Float64() * gridSize
	y := rand.Float64() * gridSize

	return &Creature{
		EntityBase: EntityBase{
			Id:     id,
			Color:  colorful.FastHappyColor().Hex(),
			X:      x,
			Y:      y,
			Radius: defaultRadius,
		},
		Name: "foo",

		dx: math.Cos(angle),
		dy: math.Sin(angle),
	}
}

func (creature *Creature) Base() *EntityBase {
	return &creature.EntityBase
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
	EntityBase
}

func RandomResource(id uint64) Entity {
	return &Resource{
		EntityBase: EntityBase{
			Id:     id,
			X:      rand.Float64() * gridSize,
			Y:      rand.Float64() * gridSize,
			Radius: resourceRadius,
			Color:  colorful.FastHappyColor().Hex(),
		},
	}
}

func (resource *Resource) Base() *EntityBase {
	return &resource.EntityBase
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
