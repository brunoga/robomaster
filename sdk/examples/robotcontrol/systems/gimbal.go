package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/brunoga/robomaster/sdk"
	"github.com/brunoga/robomaster/sdk/examples/robotcontrol/components"
	"github.com/brunoga/robomaster/sdk/modules/gimbal"
)

type systemGimbalEntity struct {
	ecs.BasicEntity
	*components.Speed
}

type Gimbal struct {
	entity        *systemGimbalEntity
	client        *sdk.Client
	mirrorClients []*sdk.Client
}

func NewGimbal(client *sdk.Client,
	mirrorClients []*sdk.Client) *Gimbal {
	return &Gimbal{
		&systemGimbalEntity{
			ecs.NewBasic(),
			&components.Speed{},
		},
		client,
		mirrorClients,
	}
}

func (g *Gimbal) New(world *ecs.World) {
	g.client.GimbalModule().Recenter()
	for _, mirrorClient := range g.mirrorClients {
		mirrorClient.GimbalModule().Recenter()
	}
}

func pixelsToDegrees(pixels float32, resolutionPixels, fovDegrees int) float64 {
	return (float64(pixels) * float64(fovDegrees)) / float64(resolutionPixels)
}

func pixelsToYawDegrees(pixels float32) float64 {
	return pixelsToDegrees(pixels, sdk.CameraVerticalResolutionPoints,
		sdk.CameraVerticalFOVDegrees)
}

func pixelsToPitchDegrees(pixels float32) float64 {
	return pixelsToDegrees(pixels, sdk.CameraHorizontalResolutionPoints,
		sdk.CameraHorizontalFOVDegrees)
}

func (g *Gimbal) Update(dt float32) {
	mouseXDelta := engo.Input.Axis("MouseXAxis").Value()
	mouseYawAngle := pixelsToYawDegrees(mouseXDelta)

	mouseYDelta := engo.Input.Axis("MouseYAxis").Value()
	mousePitchAngle := pixelsToPitchDegrees(mouseYDelta)

	if g.entity.SpeedX != mouseYawAngle ||
		g.entity.SpeedY != mousePitchAngle {
		g.client.GimbalModule().SetSpeed(gimbal.NewSpeed(-mousePitchAngle*30,
			mouseYawAngle*30), true)
		for _, mirrorClient := range g.mirrorClients {
			mirrorClient.GimbalModule().SetSpeed(gimbal.NewSpeed(
				-mousePitchAngle*30, mouseYawAngle*30), true)
		}

		g.entity.SpeedX = mouseYawAngle
		g.entity.SpeedY = mousePitchAngle
	}
}

func (g *Gimbal) Remove(e ecs.BasicEntity) {
	// Do nothing.
}

func (g *Gimbal) Priority() int {
	return 10
}
