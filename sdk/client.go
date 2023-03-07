package sdk

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/brunoga/robomaster/sdk/internal/text/modules/event"
	"github.com/brunoga/robomaster/sdk/internal/text/modules/push"
	"github.com/brunoga/robomaster/sdk/modules/robot"

	"github.com/brunoga/robomaster/sdk/internal/text/modules/armor"
	"github.com/brunoga/robomaster/sdk/internal/text/modules/blaster"
	"github.com/brunoga/robomaster/sdk/internal/text/modules/chassis"
	"github.com/brunoga/robomaster/sdk/internal/text/modules/control"
	"github.com/brunoga/robomaster/sdk/internal/text/modules/finder"
	"github.com/brunoga/robomaster/sdk/internal/text/modules/gimbal"
	"github.com/brunoga/robomaster/sdk/internal/text/modules/sound"
	"github.com/brunoga/robomaster/sdk/internal/text/modules/video"
	"github.com/brunoga/robomaster/sdk/support/logger"

	robotint "github.com/brunoga/robomaster/sdk/internal/text/modules/robot"
)

// Client enables controlling a RoboMaster robot through the plain-text SDK
// API (https://robomaster-dev.readthedocs.io/en/latest/).
type Client struct {
	logger *logger.Logger

	finderModule  *finder.Finder
	controlModule *control.Control

	pushModule  *push.Push
	eventModule *event.Event

	robotModule   *robotint.Robot
	gimbalModule  *gimbal.Gimbal
	chassisModule *chassis.Chassis
	videoModule   *video.Video
	armorModule   *armor.Armor
	blasterModule *blaster.Blaster
	soundModule   *sound.Sound
}

// NewClient returns a new client instance associated with the given ip. If ip
// is nil, the Client will try to detect a robot broadcasting its ip in the
// network.
func NewClient(ip net.IP) (*Client, error) {
	// Setup logging.
	l := logger.New(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	// Initialize all modules.
	finderModule := finder.New(l)
	if ip != nil {
		finderModule.SetIP(ip)
	}
	controlModule, err := control.New(finderModule, l)
	if err != nil {
		return nil, fmt.Errorf("error creating control module: %w", err)
	}
	pushModule, err := push.NewPush(controlModule)
	if err != nil {
		return nil, fmt.Errorf("error creating push module: %w", err)
	}
	eventModule, err := event.NewEvent(controlModule)
	if err != nil {
		return nil, fmt.Errorf("error creating event module: %w", err)
	}
	robotModule := robotint.New(controlModule)
	gimbalModule := gimbal.New(controlModule, pushModule)
	chassisModule := chassis.New(controlModule, pushModule)
	armorModule := armor.New(controlModule, eventModule)
	blasterModule := blaster.New(controlModule)
	soundModule := sound.New(eventModule)
	videoModule, err := video.New(controlModule)
	if err != nil {
		panic(err)
	}

	return &Client{
		l,
		finderModule,
		controlModule,
		pushModule,
		eventModule,
		robotModule,
		gimbalModule,
		chassisModule,
		videoModule,
		armorModule,
		blasterModule,
		soundModule,
	}, nil
}

// NewClientUSB creates a Client that tries to connect to the default USB
// connection ip.
func NewClientUSB() (*Client, error) {
	return NewClient(net.ParseIP("192.168.42.2"))
}

// NewClientWifiDirect creates a Client that tries to connect to the default
// WiFi Direct connection ip.
func NewClientWifiDirect() (*Client, error) {
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
func (c *Client) RobotModule() robot.Robot {
	return c.robotModule
}

// GimbalModule returns a pointer to the associated Gimbal module. Used for
// doing gimbal-related operations.
func (c *Client) GimbalModule() *gimbal.Gimbal {
	return c.gimbalModule
}

// ChassisModule returns a pointer to the associated Chassis module. Used for
// doing chassis-related operations.
func (c *Client) ChassisModule() *chassis.Chassis {
	return c.chassisModule
}

// VideoModule returns a pointer to the associated Video module. Used for
// doing video-related operations.
func (c *Client) VideoModule() *video.Video {
	return c.videoModule
}

// ArmorModule returns a pointer to the associated Armor module. Used for
// setting/getting hit sensitivity and detecting hits.
func (c *Client) ArmorModule() *armor.Armor {
	return c.armorModule
}

// BlasterModule returns a pointer to the associated Video module. Used for
// firing beads.
func (c *Client) BlasterModule() *blaster.Blaster {
	return c.blasterModule
}

// SoundModule returns a pointer to the associated Video module. Used for
// detecting applause (clapping) sounds.
func (c *Client) SoundModule() *sound.Sound {
	return c.soundModule
}
