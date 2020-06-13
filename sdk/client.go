package sdk

import (
	"fmt"
	"github.com/brunoga/robomaster/sdk/modules/gimbal"
	"github.com/brunoga/robomaster/sdk/modules/push"
	"github.com/brunoga/robomaster/sdk/modules/robot"
	"net"

	"github.com/brunoga/robomaster/sdk/modules"
)

// Client enables controlling a RoboMaster robot through the plain-text SDK
// API (https://robomaster-dev.readthedocs.io/en/latest/).
type Client struct {
	finderModule  *modules.Finder
	controlModule *modules.Control
	eventModule   *modules.Event
	pushModule    *push.Push

	robotModule  *robot.Robot
	gimbalModule *gimbal.Gimbal
	videoModule  *modules.Video
}

// NewClient returns a new client instance associated with the given ip. If ip
// is nil, the Client will try to detect a robot broadcasting its ip in the
// network.
func NewClient(ip net.IP) *Client {
	finderModule := modules.NewFinder()
	if ip != nil {
		finderModule.SetIP(ip)
	}

	controlModule := modules.NewControl(finderModule, false)
	eventModule := modules.NewEvent(controlModule)
	pushModule := push.NewPush(controlModule)
	robotModule := robot.NewRobot(controlModule)
	gimbalModule := gimbal.NewGimbal(controlModule, pushModule)
	videoModule := modules.NewVideo(controlModule)

	return &Client{
		finderModule,
		controlModule,
		eventModule,
		pushModule,
		robotModule,
		gimbalModule,
		videoModule,
	}
}

// NewClientUSB creates a Client that tries to connect to the default USB
// connection ip.
func NewClientUSB() *Client {
	return NewClient(net.ParseIP("192.168.42.2"))
}

// NewClientWifiDirect creates a Client that tries to connect to the default
// WiFi Direct connection ip.
func NewClientWifiDirect() *Client {
	return NewClient(net.ParseIP("192.168.2.1"))
}

// Open opens the Client connection to the robot and enters SDK mode. Returns a
// nil error on success and a non-nil error on failure.
func (c *Client) Open() error {
	err := c.controlModule.Open()
	if err != nil {
		return fmt.Errorf("error starting control module: %w", err)
	}

	// Enter SDK mode.
	err = c.controlModule.SendDataExpectOk("command;")
	if err != nil {
		return fmt.Errorf("error entering to sdk mode: %w", err)
	}

	return nil
}

// Close exits SDk mode and closes the Client connection to the robot. Returns
// a nil error on success and a non-nil error on failure.
func (c *Client) Close() error {
	// Leave SDK mode.
	err := c.controlModule.SendDataExpectOk("quit;")
	if err != nil {
		return fmt.Errorf("error leaving sdk mode: %w", err)
	}

	err = c.controlModule.Close()
	if err != nil {
		return fmt.Errorf("error stopping control module")
	}

	return nil
}

// RobotModule returns a pointer to the associated Robot module. Used for
// doing generic robot-related operations.
func (c *Client) RobotModule() *robot.Robot {
	return c.robotModule
}

// GimbalModule returns a pointer to the associated Gimbal module. Used for
// doing gimbal-related operations.
func (c *Client) GimbalModule() *gimbal.Gimbal {
	return c.gimbalModule
}

// VideoModule returns a pointer to the associated Video module. Used for
// doing video-related operations.
func (c *Client) VideoModule() *modules.Video {
	return c.videoModule
}
