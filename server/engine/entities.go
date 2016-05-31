package engine

import (
	"math"
	"math/rand"
	"strings"

	"github.com/egnwd/outgain/server/protocol"
	"github.com/lucasb-eyer/go-colorful"
)

const defaultRadius float64 = 0.5
const resourceRadius float64 = 0.1
const resourceVolume float64 = 1

type Entity interface {
	Tick(dt float64)
	Serialize() protocol.Entity
	Base() *EntityBase
	Volume() float64
}

type EntityBase struct {
	ID     uint64
	Color  string
	X      float64
	Y      float64
	Radius float64

	// These are used to work around that the list of entities should
	// not be modified, nor should the entities' radii and coordinates
	// while the collision algorithm is running
	//
	// Instead we track modifications which need to be performed here,
	// and perform them afterwards
	dying      bool
	nextRadius float64
}

// Left gets the X coordinate of the left hand side
func (entity *EntityBase) Left() float64 {
	return entity.X - entity.Radius
}

// Right gets the X coordinate of the right hand side
func (entity *EntityBase) Right() float64 {
	return entity.X + entity.Radius
}

// Top gets the Y coordinate of the top side
func (entity *EntityBase) Top() float64 {
	return entity.Y - entity.Radius
}

// Bottom gets the Y coordinate of the bottom side
func (entity *EntityBase) Bottom() float64 {
	return entity.Y + entity.Radius
}

// EntityList is shorthand for a slice of Entitys
type EntityList []Entity

func (list EntityList) Len() int {
	return len(list)
}

func (list EntityList) Less(i, j int) bool {
	return list[i].Base().Left() < list[j].Base().Left()
}

func (list EntityList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

// Tick every entity of the list
func (list EntityList) Tick(dt float64) {
	for _, entity := range list {
		entity.Tick(dt)
	}

	// Ticking could have moved some entities, so sort the list again to
	// maintain the invariant
	list.Sort()

}

func (list EntityList) Filter(filter func(Entity) bool) EntityList {
	count := list.Len()
	for i := 0; i < count; i++ {
		if !filter(list[i]) {
			list.Swap(i, count-1)
			count--
		}
	}
	return list[:count]
}

func (list EntityList) Insert(entity Entity) EntityList {
	result := append(list, entity)

	for i := result.Len() - 1; i > 0 && result.Less(i, i-1); i-- {
		result.Swap(i-1, i)
	}

	return result
}

func (list EntityList) Sort() {
	for i := 1; i < list.Len(); i++ {
		for j := i; j > 0 && list.Less(j, j-1); j-- {
			list.Swap(j-1, j)
		}
	}
}

type Creature struct {
	EntityBase

	Name   string
	Sprite string

	dx float64
	dy float64
}

func RandomCreature(id uint64, name string) Entity {
	angle := rand.Float64() * 2 * math.Pi
	x := rand.Float64() * gridSize
	y := rand.Float64() * gridSize
	color := colorful.FastHappyColor().Hex()

	return &Creature{
		EntityBase: EntityBase{
			ID:     id,
			Color:  color,
			X:      x,
			Y:      y,
			Radius: defaultRadius,
		},
		Name:   name,
		Sprite: "/images/creature-" + strings.TrimPrefix(color, "#") + ".png",

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
		ID:     creature.ID,
		Name:   &creature.Name,
		Sprite: &creature.Sprite,
		Color:  creature.Color,
		X:      creature.X,
		Y:      creature.Y,
		Radius: creature.Radius,
	}
}

func (creature *Creature) Volume() float64 {
	return creature.nextRadius * creature.nextRadius
}

type Resource struct {
	EntityBase
}

func RandomResource(id uint64, _ string) Entity {
	return &Resource{
		EntityBase: EntityBase{
			ID:     id,
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
		ID:     resource.ID,
		Name:   nil,
		Sprite: nil,
		Color:  resource.Color,
		X:      resource.X,
		Y:      resource.Y,
		Radius: resource.Radius,
	}
}

func (resource *Resource) Volume() float64 {
	return resourceVolume
}
