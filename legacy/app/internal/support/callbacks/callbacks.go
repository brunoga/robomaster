package callbacks

import (
	"fmt"
	"reflect"
	"sync"
)

// Key represents a primary key for the callback map. It will usually be either
// an unity.EventType or a dji.Key.
type Key uint64

// Tag represents a secondary key for the callback map. It is used to track
// callbacks so we are able to remove a continuous callback (or, internally, to
// auto-remove a single-shot callback).
type Tag uint64

// KeyFunc is the prototype for functions that should run when the first
// callback for a key is added and the last one is removed.
type KeyFunc func(key Key) error

type data struct {
	callback interface{}
	once     bool
}

// Callbacks holds a map of callbacks (continuous or single-shot).
type Callbacks struct {
	name      string
	firstFunc KeyFunc
	lastFunc  KeyFunc

	m           sync.Mutex
	callbackMap map[Key]map[Tag]*data
	nextTag     uint64
}

// New returns a new Callbacks instance with the given name, firstFunc and
// lastFunc. The firstFunc function is run before the first callback is added
// to the map and the lastFunc is run after the last one is removed. Either
// might be nil which means do not run.
func New(name string, firstFunc KeyFunc, lastFunc KeyFunc) *Callbacks {
	return &Callbacks{
		name,
		firstFunc,
		lastFunc,
		sync.Mutex{},
		make(map[Key]map[Tag]*data),
		0,
	}
}

// AddSingleShot adds a new callback that will only fire once. The callback will
// be automatically removed from the map when Callback() returns it. Returns a
// nil error on success and a non-nil error on failure.
func (c *Callbacks) AddSingleShot(key Key, callback interface{}) (Tag, error) {
	return c.add(key, callback, true)
}

// AddContinous adds a new callback that will fire continuously. When needed it
// has to be explicitly removed. Returns an associated Tag and nil error on
// success and Tag(0) and a non-nil error on failure.
func (c *Callbacks) AddContinuous(key Key, callback interface{}) (Tag, error) {
	return c.add(key, callback, false)
}

func (c *Callbacks) add(key Key, callback interface{}, once bool) (Tag, error) {
	if callback == nil {
		return 0, fmt.Errorf("%s : callback must not be nil", c.name)
	}

	reflectionType := reflect.TypeOf(callback)
	if reflectionType.Kind() != reflect.Func {
		return 0, fmt.Errorf("%s : callback must be a function", c.name)
	}

	c.m.Lock()
	defer c.m.Unlock()

	tagMap, ok := c.callbackMap[key]
	if !ok {
		if c.firstFunc != nil {
			err := c.firstFunc(key)
			if err != nil {
				return 0, err
			}
		}
		tagMap = make(map[Tag]*data)
		c.callbackMap[key] = tagMap
	}

	// Ensures that the first Tag will be 1.
	c.nextTag++

	tag := Tag(c.nextTag)

	tagMap[tag] = &data{
		callback,
		once,
	}

	return tag, nil
}

// Remove removes a callback that matches the given key or tag. Returns a nil
// error on success and a non-nil error in failure. Returns a nil error on
// success and a non-nil error on failure.
func (c *Callbacks) Remove(key Key, tag Tag) error {
	c.m.Lock()
	defer c.m.Unlock()

	return c.remove(key, tag, false)
}

func (c *Callbacks) remove(key Key, tag Tag, allowOnce bool) error {
	if tag == 0 {
		return fmt.Errorf("%s : invalid tag", c.name)
	}

	tagMap, ok := c.callbackMap[key]
	if !ok {
		return fmt.Errorf("%s : key not found", c.name)
	}

	d, ok := tagMap[tag]
	if !ok {
		return fmt.Errorf("%s : tag not found for given key", c.name)
	}

	if d.once && !allowOnce {
		return fmt.Errorf("%s : can not remove single-use callback",
			c.name)
	}

	delete(tagMap, tag)

	var err error = nil
	if len(tagMap) == 0 {
		delete(c.callbackMap, key)

		if c.lastFunc != nil {
			err = c.lastFunc(key)
		}
	}

	return err
}

// Callback returns the callback associated with Key and Tag. If the callback is
// one-shot, it will also automatically remove it. Returns the callback (as an
// interface{}) and a nil error on success and nil and a non-nil error on
// failure.
func (c *Callbacks) Callback(key Key, tag Tag) (interface{}, error) {
	c.m.Lock()
	defer c.m.Unlock()

	if tag == 0 {
		return nil, fmt.Errorf("%s : invalid tag", c.name)
	}

	d, err := c.getData(key, tag)
	if err != nil {
		return nil, err
	}

	if d.once {
		err = c.remove(key, tag, true)
		if err != nil {
			return nil, err
		}
	}

	return d.callback, nil
}

func (c *Callbacks) getData(key Key, tag Tag) (*data, error) {
	tagMap, ok := c.callbackMap[key]
	if !ok {
		return nil, fmt.Errorf("%s : key not found", c.name)
	}

	d, ok := tagMap[tag]
	if !ok {
		return nil, fmt.Errorf("%s : tag not found for given key",
			c.name)
	}

	return d, nil
}

// CallbacksForKey returns all callbacks associated with the given key. Returns
// a slice of callbacks (as interfaces{}) and a nil error on success and nil and
// a non-nil error on failure.
func (c *Callbacks) CallbacksForKey(key Key) ([]interface{}, error) {
	c.m.Lock()
	defer c.m.Unlock()

	tagMap, ok := c.callbackMap[key]
	if !ok {
		return nil, fmt.Errorf("%s : key not found", c.name)
	}

	cbs := make([]interface{}, len(tagMap))

	i := 0
	for tag, d := range tagMap {
		if d.once {
			err := c.remove(key, tag, true)
			if err != nil {
				return nil, err
			}
		}

		cbs[i] = d.callback
		i++
	}

	return cbs, nil
}
