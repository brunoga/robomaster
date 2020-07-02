package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster/sdk"
)

type Blaster struct {
	client        *sdk.Client
	mirrorClients []*sdk.Client
}

func NewBlaster(client *sdk.Client, mirrorClients []*sdk.Client) *Blaster {
	return &Blaster{
		client,
		mirrorClients,
	}
}

func (b *Blaster) New(world *ecs.World) {
	b.client.BlasterModule().SetNumBeads(1)
	for _, mirrorClient := range b.mirrorClients {
		mirrorClient.BlasterModule().SetNumBeads(1)
	}
}

func (b *Blaster) Update(dt float32) {
	if engo.Input.Mouse.Action == engo.Press &&
		engo.Input.Mouse.Button == engo.MouseButtonLeft {
		b.client.BlasterModule().Fire(true)
		for _, mirrorClient := range b.mirrorClients {
			mirrorClient.BlasterModule().Fire(true)
		}
	}
}

func (b *Blaster) Remove(e ecs.BasicEntity) {
	// Do nothing.
}

func (b *Blaster) Priority() int {
	return 11
}
