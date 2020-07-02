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
	*components.Position
}

type Gimbal struct {
	entity *systemGimbalEntity
	client *sdk.Client
}

func NewGimbal(client *sdk.Client) *Gimbal {
	return &Gimbal{
		&systemGimbalEntity{
			ecs.NewBasic(),
			&components.Position{},
		},
		client,
	}
}

func (g *Gimbal) New(world *ecs.World) {
	g.client.GimbalModule().Recenter()
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

	if g.entity.PositionX > mouseYawAngle ||
		g.entity.PositionY != mousePitchAngle {
		//g.client.GimbalModule().MoveRelative(
		//	gimbal.NewPosition(-mousePitchAngle, mouseYawAngle),
		//	gimbal.NewSpeed(mousePitchAngle*30, mouseYawAngle*30),
		//	true)
		g.client.GimbalModule().SetSpeed(gimbal.NewSpeed(-mousePitchAngle*30, mouseYawAngle*30), true)
		g.entity.PositionX = mouseYawAngle
		g.entity.PositionY = mousePitchAngle
	}
}

func (g *Gimbal) Remove(e ecs.BasicEntity) {
	// Do nothing.
}

func (g *Gimbal) Priority() int {
	return 10
}
