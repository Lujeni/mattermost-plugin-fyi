package main

import (
	"fmt"

	//	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
)

type configuration struct {
	Tags          string
	GrafanaURL    string
	GrafanaAPIKey string
}

func (c *configuration) Clone() *configuration {
	var clone = *c
	return &clone
}

// IsValid checks if all needed fields are set.
func (c *configuration) IsValid() error {
	if c.GrafanaAPIKey == "" {
		return fmt.Errorf("must have a grafana api key")
	}

	if c.GrafanaURL == "" {
		return fmt.Errorf("must have a grafana url")
	}

	return nil
}

func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		panic("setConfiguration called with the existing configuration")
	}

	p.configuration = configuration
}

func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin configuration")
	}

	p.setConfiguration(configuration)

	return nil
}
