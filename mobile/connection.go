package mobile

import (
	"github.com/brunoga/robomaster/module/connection"
)

// Camera allows controlling the robot's connection.
type Connection struct {
	c *connection.Connection
}

// SignalQualityLevel returns the current signal quality level (0 to 60).
func (c *Connection) SignalQualityLevel() int8 {
	return int8(c.c.SignalQualityLevel())
}

// SignalQualityBars returns the current signal quality bars (1 to 4).
func (c *Connection) SignalQualityBars() int8 {
	return int8(c.c.SignalQualityBars())
}
