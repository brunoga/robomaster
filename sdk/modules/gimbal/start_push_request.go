package gimbal

import (
	"github.com/brunoga/robomaster/sdk/modules/push"
)

type StartPushRequest struct {
	PushAttribute PushAttribute
	PushHandler   push.Handler
}
