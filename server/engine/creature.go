package engine

import (
	"log"
	"math"
	"math/rand"
	"strings"

	"github.com/egnwd/outgain/server/config"
	"github.com/egnwd/outgain/server/guest"
	"github.com/egnwd/outgain/server/protocol"
	"github.com/egnwd/outgain/server/runner"

	"github.com/lucasb-eyer/go-colorful"
)

const speedFactor = 4

type Creature struct {
	EntityBase

	Guest  *guest.Guest
	Sprite string

	runner *runner.RunnerClient
}

func (creature *Creature) incrementScore(eaten Entity) {
	creature.Guest.AddGains(1)
}

func (creature *Creature) decrementScore() {
	creature.Guest.LoseGains(1)
}

func NewCreature(guest *guest.Guest, config *config.Config) (builderFunc, error) {
	x := rand.Float64() * gridSize
	y := rand.Float64() * gridSize
	random_light := colorful.Hcl(rand.Float64()*360.0, rand.Float64(), 0.6+rand.Float64()*0.4)
	color := random_light.Hex()

	client, err := runner.StartRunner(config)
	if err != nil {
		return nil, err
	}

	err = client.Load(guest.Source)
	if err != nil {
		return nil, err
	}

	return func(id uint64) Entity {
		return &Creature{
			EntityBase: EntityBase{
				ID:     id,
				Color:  color,
				X:      x,
				Y:      y,
				Radius: defaultRadius,
			},
			Guest: guest,

			Sprite: "/images/creature-" + strings.TrimPrefix(color, "#") + ".svg",
			runner: client,
		}
	}, nil
}

func (creature *Creature) GetName() string {
	return creature.Guest.GetName()
}

func (creature *Creature) GetGains() int {
	return creature.Guest.GetGains()
}

func (creature *Creature) Base() *EntityBase {
	return &creature.EntityBase
}

func (creature *Creature) Tick(state protocol.WorldState, dt float64) {
	speed, err := creature.runner.Tick(creature.Serialize(), state)
	if err != nil {
		log.Printf("Creature %s: %v", creature.GetName(), err)
		creature.dying = true
		return
	}

	norm := math.Sqrt(speed.Dx*speed.Dx + speed.Dy*speed.Dy)
	if norm > 1 {
		speed.Dx /= norm
		speed.Dy /= norm
	}

	creature.X += speed.Dx * dt * speedFactor
	creature.Y += speed.Dy * dt * speedFactor

	if creature.X-creature.Radius < 0 {
		creature.X = creature.Radius
	}
	if creature.X+creature.Radius > gridSize {
		creature.X = gridSize - creature.Radius
	}
	if creature.Y-creature.Radius < 0 {
		creature.Y = creature.Radius
	}
	if creature.Y+creature.Radius > gridSize {
		creature.Y = gridSize - creature.Radius
	}
}

func (creature *Creature) Serialize() protocol.Entity {
	return protocol.Entity{
		ID:         creature.ID,
		Name:       &creature.Guest.Name,
		Sprite:     &creature.Sprite,
		Color:      creature.Color,
		X:          creature.X,
		Y:          creature.Y,
		Radius:     creature.Radius,
		EntityType: creatureEnum,
	}
}

func (creature *Creature) BonusFactor() float64 {
	return 1
}

func (creature *Creature) Close() {
	creature.runner.Close()
}

func (creature *Creature) IsUser() bool {
	return creature.Guest.IsUser()
}
