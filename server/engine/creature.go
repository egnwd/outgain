package engine

import (
	"io/ioutil"
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

type Creature struct {
	EntityBase

	Guest  *guest.Guest
	Sprite string

	runner *runner.RunnerClient
}

func (creature *Creature) incrementScore(eaten Entity) {
	creature.Guest.AddResources(1)
}

func NewCreature(guest *guest.Guest, config *config.Config) func(id uint64) Entity {
	return func(id uint64) Entity {
		x := rand.Float64() * gridSize
		y := rand.Float64() * gridSize
		color := colorful.FastHappyColor().Hex()

		source, err := ioutil.ReadFile(config.DefaultAI)
		if err != nil {
			log.Fatalln(err)
		}

		client, err := runner.StartRunner(config, string(source))
		if err != nil {
			log.Fatalln(err)
		}

		return &Creature{
			EntityBase: EntityBase{
				ID:     id,
				Color:  color,
				X:      x,
				Y:      y,
				Radius: defaultRadius,
			},
			Guest:  guest,
			Sprite: "/images/creature-" + strings.TrimPrefix(color, "#") + ".svg",
			runner: client,
		}
	}
}

func (creature *Creature) GetName() string {
	return creature.Guest.GetName()
}

func (creature *Creature) GetResources() int {
	return creature.Guest.GetResources()
}

func (creature *Creature) Base() *EntityBase {
	return &creature.EntityBase
}

func (creature *Creature) Tick(state protocol.WorldState, dt float64) {
	speed, err := creature.runner.Tick(creature.Serialize(), state)
	if err != nil {
		log.Fatalln(err)
	}

	norm := math.Sqrt(speed.Dx*speed.Dx + speed.Dy*speed.Dy)
	if norm > 1 {
		speed.Dx /= norm
		speed.Dy /= norm
	}

	creature.X += speed.Dx * dt
	creature.Y += speed.Dy * dt

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
		ID:     creature.ID,
		Name:   &creature.Guest.Name,
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

func (creature *Creature) Close() {
	creature.runner.Close()
}
