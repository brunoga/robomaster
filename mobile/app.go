package mobile

import "github.com/brunoga/robomaster/unitybridge/support"

// App holds the ID for a specific app.
type App struct {
	id int64
}

// NewApp creates a new App instance with the given ID. If the ID is 0, a new
// ID will be generated.
func NewApp(id int64) (*App, error) {
	if id == 0 {
		newID, err := support.GenerateAppID()
		if err != nil {
			return nil, err
		}
		id = int64(newID)
	}

	return &App{
		id: id,
	}, nil
}

// ID returns the ID for the App instance.
func (a *App) ID() int64 {
	return a.id
}
