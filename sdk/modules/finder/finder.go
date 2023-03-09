package finder

import (
	"time"
)

type Finder interface {
	Find(time.Duration, Func) error
}
