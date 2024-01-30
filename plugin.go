package traefikaggregator

import (
	"context"
	"encoding/json"
	"fmt"

	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/application"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/config"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/log"
)

// We need this defined here, because the plugin loader does not support importing from other packages
type Config struct {
	PollInterval string                   `json:"pollInterval,omitempty"`
	Instances    []config.TraefikInstance `json:"instances,omitempty"`
}

type Provider struct {
	app application.App
}

func CreateConfig() *config.Config {
	return &config.Config{}
}

func New(ctx context.Context, config *config.Config, name string) (*Provider, error) {
	log.Info("New Traefik Aggregator plugin with name %s and config %+v", name, config)

	app, err := application.NewApp(config)

	if err != nil {
		return nil, fmt.Errorf("error while creating app: %w", err)
	}

	return &Provider{
		app: app,
	}, nil
}

func (p *Provider) Init() error {
	log.Info("Initializing Traefik Aggregator plugin")

	return nil
}

func (p *Provider) Provide(cfgChan chan<- json.Marshaler) error {
	ctx, _ := context.WithCancel(context.Background())

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Recovered from panic: %v", err)
			}
		}()

		p.app.Run(ctx, cfgChan)
	}()

	return nil
}

func (p *Provider) Stop() error {
	return nil
}
