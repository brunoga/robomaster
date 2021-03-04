package finder

import (
	"time"

	"github.com/brunoga/robomaster/sdk/robot"
)

// Filter is a map used to pass filtering information to a Finder
// implementation. Any unknown keys are ignored. A nil filter is valid and
// means all robots.
type Filter map[string]interface{}

// Finder is the interface for detecting and enumerating robots.
type Finder interface {
	// Find looks for broadcasts from robots that satisfy the filter criteria for
	// the specified duration. Any robots not detected within the given duration
	// will be ignored.
	Find(filter Filter, duration time.Duration)

	// NumRobots returns the number of detected robots (so far, in case the
	// duration did not expire yet) that satisfy the given filter.
	NumRobots(filter Filter) int

	// Robots returns all robots that match the given filter or nil in case no
	// robot satisfies the filter (or there are less than n robots that do). The
	// given filter may include an entry with key:n and value:int to return the
	// single nth robot from the robot list after it is filtered by any other
	// filter keys (if one exists).
	Robots(filter Filter) []robot.Robot
}
