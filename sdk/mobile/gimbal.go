package mobile

import "github.com/brunoga/robomaster/sdk/modules/gimbal"

type Gimbal struct {
	g *gimbal.Gimbal
}

type GimbalSpeed struct {
	Pitch float64
	Yaw   float64
}

type GimbalPosition struct {
	Pitch float64
	Yaw   float64
}

type GimbalAttitude struct {
	Pitch float64
	Yaw   float64
}

func (g *Gimbal) SetSpeed(speed *GimbalSpeed, async bool) error {
	return g.g.SetSpeed(gimbal.NewSpeed(speed.Pitch, speed.Yaw), async)
}

func (g *Gimbal) MoveRelative(position *GimbalPosition, speed *GimbalSpeed,
	async bool) error {
	return g.g.MoveRelative(gimbal.NewPosition(position.Pitch, position.Yaw),
		gimbal.NewSpeed(speed.Pitch, speed.Yaw), async)
}

func (g *Gimbal) MoveAbsolute(position *GimbalPosition, speed *GimbalSpeed,
	async bool) error {
	return g.g.MoveAbsolute(gimbal.NewPosition(position.Pitch, position.Yaw),
		gimbal.NewSpeed(speed.Pitch, speed.Yaw), async)
}

func (g *Gimbal) Suspend() error {
	return g.g.Suspend()
}

func (g *Gimbal) Resume() error {
	return g.g.Resume()
}

func (g *Gimbal) Recenter() error {
	return g.g.Recenter()
}

func (g *Gimbal) GetAttitude() (*GimbalAttitude, error) {
	attitude, err := g.g.GetAttitude()
	if err != nil {
		return nil, err
	}

	return &GimbalAttitude{
		Pitch: attitude.Pitch(),
		Yaw:   attitude.Yaw(),
	}, nil
}

func (g *Gimbal) StartPush(attr int,
	pushHandler NotificationHandler) (int, error) {
	return g.g.StartPush(gimbal.PushAttribute(attr), pushHandler.Handle)
}

func (g *Gimbal) StopPush(attr int, token int) error {
	return g.g.StopPush(gimbal.PushAttribute(attr), token)
}
