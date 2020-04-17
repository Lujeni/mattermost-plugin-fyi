package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

type Plugin struct {
	plugin.MattermostPlugin
	configurationLock sync.RWMutex
	configuration     *configuration
}

func (p *Plugin) OnActivate() error {
	config := p.getConfiguration()

	if err := config.IsValid(); err != nil {
		return err
	}
	if err := p.API.RegisterCommand(getCommand()); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Unable to register command: %v", getCommand()))
	}

	return nil
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}
