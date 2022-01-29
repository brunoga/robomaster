package components

import (
	"git.bug-br.org.br/bga/robomasters1/app"
)

type Robomaster struct {
	a *app.App
}

func NewRobomaster(ssId, password string, appID uint64) (*Robomaster, error) {
	var a *app.App
	var err error
	if appID != 0 {
		a, err = app.NewWithAppID("US", ssId, password, "", appID)
	} else {
		a, err = app.New("US", ssId, password, "")
	}
	if err != nil {
		return nil, err
	}

	return &Robomaster{
		a,
	}, nil
}

func (r *Robomaster) App() *app.App {
	return r.a
}
