package engine

import (
	"log"
	"math/rand"
	"strings"

	"github.com/egnwd/outgain/server/config"
	"github.com/egnwd/outgain/server/protocol"
	"github.com/egnwd/outgain/server/runner"

	"github.com/lucasb-eyer/go-colorful"
)

type Creature struct {
	EntityBase

	Name   string
	Sprite string

	runner *runner.RunnerClient
}

const creatureAI string = `
angle = rand * 2 * Math::PI
@dx = Math::cos(angle)
@dy = Math::sin(angle)

def run
    angle = rand * Math::PI / 2 - Math::PI / 4

    c = Math::cos(angle)
    s = Math::sin(angle)

    dx = c * @dx - s * @dy
    dy = s * @dx + c * @dy

    @dx = dx
    @dy = dy

    move(@dx, @dy)
end
`

func RandomCreature(config *config.Config) func(id uint64) Entity {
	return func(id uint64) Entity {
		x := rand.Float64() * gridSize
		y := rand.Float64() * gridSize
		color := colorful.FastHappyColor().Hex()

		client, err := runner.StartRunner(config, creatureAI)
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
			Name:   "foo",
			Sprite: "/images/creature-" + strings.TrimPrefix(color, "#") + ".png",
			runner: client,
		}
	}
}

func (creature *Creature) Base() *EntityBase {
	return &creature.EntityBase
}

func (creature *Creature) Tick(dt float64) {
	movement, err := creature.runner.Tick(protocol.WorldState{})
	if err != nil {
		log.Fatalln(err)
	}

	creature.X += movement.Dx * dt
	creature.Y += movement.Dy * dt

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

func (creature *Creature) Close() {
	creature.runner.Close()
}
