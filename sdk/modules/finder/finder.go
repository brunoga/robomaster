package finder

import (
	"time"

	"github.com/brunoga/robomaster/sdk/modules/robot"
)

// Filter is a map used to pass filtering information to a Finder
// implementation m ethod. Any unknown keys are ignored. A nil filter
// is valid and means all detected robots.
type Filter map[string]interface{}

// Finder is the interface for detecting and enumerating robots.
type Finder interface {
	// Find looks for broadcasts from all robots that support the specific
	// protocol implementation and satisfy the given filter for the specified
	// duration. Any robots not detected within the given duration will be
	// ignored.
	Find(filter Filter, duration time.Duration)

	// NumRobots returns the number of detected robots (so far, in case the
	// duration did not expire yet).
	NumRobots() int

	// Robot returns the 0-based nth detected Robot (so far, in case the duration
	// did not expire yet). Returns nil if there are less than n+1 detected
	// Robots.
	Robot(n int) robot.Robot
}
