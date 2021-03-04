package internal

import "github.com/brunoga/robomaster/sdk/modules/finder"

// GetFilterParameter returns the value (as an interface{}) in the given filter
// associated with the given key. If key is not found, returns nil.
func GetFilterParameter(key string, filter finder.Filter) interface{} {
	v, ok := filter[key]
	if !ok {
		return nil
	}

	return v
}
