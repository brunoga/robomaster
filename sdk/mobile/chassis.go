package mobile

import (
	"github.com/brunoga/robomaster/sdk/modules/chassis"
)

type ChassisSpeed struct {
	X float64
	Y float64
	Z float64
}

type ChassisWheelSpeed struct {
	W1 float64
	W2 float64
	W3 float64
	W4 float64
}

type ChassisPosition struct {
	X float64
	Y float64
	Z float64
}

type ChassisSpeeds struct {
	Speed      *ChassisSpeed
	WheelSpeed *ChassisWheelSpeed
}

type ChassisAttitude struct {
	Pitch float64
	Roll  float64
	Yaw   float64
}

type ChassisStatus struct {
	IsStatic   bool
	UpHill     bool
	DownHill   bool
	OnSlope    bool
	PickUp     bool
	Slip       bool
	ImpactX    bool
	ImpactY    bool
	ImpactZ    bool
	RollOver   bool
	HillStatic bool
}

type Chassis struct {
	c *chassis.Chassis
}

func (c *Chassis) SetSpeed(speed *ChassisSpeed, async bool) error {
	convSpeed := chassis.NewSpeed(speed.X, speed.Y, speed.Z)
	return c.c.SetSpeed(convSpeed, async)
}

func (c *Chassis) SetWheelSpeed(wheelSpeed *ChassisWheelSpeed, async bool) error {
	convWheelSpeed := chassis.NewWheelSpeed(wheelSpeed.W1,
		wheelSpeed.W2, wheelSpeed.W3, wheelSpeed.W4)
	return c.c.SetWheelSpeed(convWheelSpeed, async)
}

func (c *Chassis) MoveRelative(position *ChassisPosition, speed *ChassisSpeed,
	async bool) error {
	convPosition := chassis.NewPosition(position.X, position.Y, position.Z)
	convSpeed := chassis.NewSpeed(speed.X, speed.Y, speed.Z)
	return c.c.MoveRelative(convPosition, convSpeed, async)
}

func (c *Chassis) GetSpeed() (*ChassisSpeeds, error) {
	speed, wheelSpeed, err := c.c.GetSpeed()
	if err != nil {
		return nil, err
	}

	return &ChassisSpeeds{
		Speed: &ChassisSpeed{
			X: speed.X(),
			Y: speed.Y(),
			Z: speed.Z(),
		},
		WheelSpeed: &ChassisWheelSpeed{
			W1: wheelSpeed.W1(),
			W2: wheelSpeed.W2(),
			W3: wheelSpeed.W3(),
			W4: wheelSpeed.W4(),
		},
	}, nil
}

func (c *Chassis) GetPosition() (*ChassisPosition, error) {
	position, err := c.c.GetPosition()
	if err != nil {
		return nil, err
	}

	return &ChassisPosition{
		X: position.X(),
		Y: position.Y(),
		Z: position.Z(),
	}, nil
}

func (c *Chassis) GetAttitude() (*ChassisAttitude, error) {
	attitude, err := c.c.GetAttitude()
	if err != nil {
		return nil, err
	}

	return &ChassisAttitude{
		Pitch: attitude.Pitch(),
		Roll:  attitude.Roll(),
		Yaw:   attitude.Yaw(),
	}, nil
}

func (c *Chassis) GetStatus() (*ChassisStatus, error) {
	status, err := c.c.GetStatus()
	if err != nil {
		return nil, err
	}

	return &ChassisStatus{
		IsStatic:   status.IsStatic(),
		UpHill:     status.IsUphill(),
		DownHill:   status.IsDownhill(),
		OnSlope:    status.IsOnSlope(),
		PickUp:     status.IsPickedUp(),
		Slip:       status.IsSlipping(),
		ImpactX:    status.XImpactDetected(),
		ImpactY:    status.YImpactDetected(),
		ImpactZ:    status.ZImpactDetected(),
		RollOver:   status.IsRolledOver(),
		HillStatic: status.IsStaticOnHill(),
	}, nil
}

func (c *Chassis) StartPush(pushAttribute int,
	pushHandler NotificationHandler, frequency int) (int, error) {
	return c.c.StartPush(chassis.PushAttribute(pushAttribute), pushHandler.Handle, frequency)
}
