package gimbal

import (
	"fmt"
	"sync"
	"time"

	"git.bug-br.org.br/bga/robomasters1/app/internal"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji"
)

type Gimbal struct {
	cc *internal.CommandController

	connectionWg *sync.WaitGroup
}

func New(cc *internal.CommandController) *Gimbal {
	connectionWg := sync.WaitGroup{}

	connectionWg.Add(1)
	cc.StartListening(dji.KeyGimbalConnection,
		func(result *dji.Result, wg *sync.WaitGroup) {
			if result.Value().(bool) {
				// Enable chassis and gimbal updates.
				fmt.Println("Gimbal connection established.")
				cc.PerformAction(
					dji.KeyRobomasterOpenChassisSpeedUpdates, nil,
					nil)
				cc.PerformAction(dji.KeyGimbalOpenAttitudeUpdates, nil,
					nil)
				connectionWg.Done()
			} else {
				fmt.Println("Gimbal connection failed.")
			}

			wg.Done()
		})

	connectionWg.Add(1)
	cc.StartListening(dji.KeyRobomasterSystemConnection,
		func(result *dji.Result, wg *sync.WaitGroup) {
			if result.Value().(bool) {
				fmt.Println("System connection established.")
				connectionWg.Done()
			} else {
				fmt.Println("System connection failed.")
			}

			wg.Done()
		})

	return &Gimbal{
		cc,
		&connectionWg,
	}
}

func (g *Gimbal) ResetPosition() error {
	return g.cc.PerformAction(dji.KeyGimbalResetPosition, nil, nil)
}

func (g *Gimbal) MoveToAbsolutePosition(yawAngle, pitchAngle int,
	duration time.Duration) error {

	param := absoluteRotationParameter{
		Time: int16(duration.Milliseconds()),
	}

	if yawAngle != 0 {
		param.Pitch = 0
		param.Yaw = int16(yawAngle * 10)
		err := g.cc.PerformAction(dji.KeyGimbalAngleFrontYawRotation,
			param, nil)
		if err != nil {
			return err
		}
	}

	if pitchAngle != 0 {
		param.Pitch = int16(pitchAngle * 10)
		param.Yaw = 0
		err := g.cc.PerformAction(dji.KeyGimbalAngleFrontPitchRotation,
			param, nil)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Gimbal) WaitForConnection() {
	g.connectionWg.Wait()
}
