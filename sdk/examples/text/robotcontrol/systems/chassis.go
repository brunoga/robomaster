package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster/sdk"
	"github.com/brunoga/robomaster/sdk/examples/text/robotcontrol/components"
	"github.com/brunoga/robomaster/sdk/text/modules/chassis"
)

type systemChassisEntity struct {
	ecs.BasicEntity
	*components.Speed
}

type Chassis struct {
	entity        *systemChassisEntity
	client        *sdk.Client
	mirrorClients []*sdk.Client
}

func NewChassis(client *sdk.Client, mirrorClients []*sdk.Client) *Chassis {
	return &Chassis{
		&systemChassisEntity{
			ecs.NewBasic(),
			&components.Speed{},
		},
		client,
		mirrorClients,
	}
}

func (c *Chassis) New(world *ecs.World) {
	// Do nothing.
}

func (c *Chassis) Update(dt float32) {
	currentLeftRight := engo.Input.Axis("Left/Right").Value()
	currentForwardBackward := -engo.Input.Axis("Forward/Backward").Value()

	if c.entity.SpeedX != float64(currentForwardBackward) ||
		c.entity.SpeedY != float64(currentLeftRight) {
		c.client.ChassisModule().SetSpeed(chassis.NewSpeed(
			float64(currentForwardBackward)/2.0, float64(currentLeftRight)/2.0,
			0.0), true)
		for _, mirrorClient := range c.mirrorClients {
			mirrorClient.ChassisModule().SetSpeed(chassis.NewSpeed(
				float64(currentForwardBackward)/2.0,
				float64(currentLeftRight)/2.0, 0.0), true)
		}

		c.entity.SpeedX = float64(currentForwardBackward)
		c.entity.SpeedY = float64(currentLeftRight)
	}
}

func (c *Chassis) Remove(e ecs.BasicEntity) {
	// Do nothing.
}

func (c *Chassis) Priority() int {
	return 9
}
