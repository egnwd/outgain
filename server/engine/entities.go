package engine

import (
	"math/rand"

	"github.com/egnwd/outgain/server/protocol"
	"github.com/lucasb-eyer/go-colorful"
)

const defaultRadius float64 = 0.35
const resourceRadius float64 = 0.1
const resourceVolume float64 = 1
const spikeRadius float64 = 0.2
const spikeVolume float64 = 1

const (
	creatureEnum = iota
	resourceEnum
	spikeEnum
)

type Entity interface {
	Tick(state protocol.WorldState, dt float64)
	Serialize() protocol.Entity
	Base() *EntityBase
	Volume() float64
	Close()
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
func (list EntityList) Tick(state protocol.WorldState, dt float64) {
	for _, entity := range list {
		entity.Tick(state, dt)
	}

	// Ticking could have moved some entities, so sort the list again to
	// maintain the invariant
	list.Sort()
}

func (list EntityList) Filter(filter func(Entity) bool) EntityList {
	count := list.Len()
	for i := 0; i < count; i++ {
		if !filter(list[i]) {
			list[i].Close()
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

type Resource struct {
	EntityBase
}

func RandomResource(id uint64) Entity {
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

func (resource *Resource) Tick(state protocol.WorldState, dt float64) {
}

func (resource *Resource) Serialize() protocol.Entity {
	return protocol.Entity{
		ID:         resource.ID,
		Name:       nil,
		Sprite:     nil,
		Color:      resource.Color,
		X:          resource.X,
		Y:          resource.Y,
		Radius:     resource.Radius,
		EntityType: resourceEnum,
	}
}

func (resource *Resource) Volume() float64 {
	return resourceVolume
}

type Spike struct {
	EntityBase
}

func RandomSpike(id uint64) Entity {
	return &Spike{
		EntityBase: EntityBase{
			ID:     id,
			X:      rand.Float64() * gridSize, // Update these so that it's not on a player
			Y:      rand.Float64() * gridSize,
			Radius: spikeRadius,
			Color:  "",
		},
	}
}

func (spike *Spike) Base() *EntityBase {
	return &spike.EntityBase
}

func (spike *Spike) Tick(dt float64) {
}

func (spike *Spike) Serialize() protocol.Entity {
	return protocol.Entity{
		ID:         spike.ID,
		Name:       nil,
		Sprite:     nil,
		Color:      spike.Color,
		X:          spike.X,
		Y:          spike.Y,
		Radius:     spike.Radius,
		EntityType: spikeEnum,
	}
}

func (spike *Spike) Volume() float64 {
	return spikeVolume
}

func (resource *Resource) Close() {
}
