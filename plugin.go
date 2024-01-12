package traefikaggregator

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/traefik/genconf/dynamic"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/aggregator"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/config"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/log"
)

type Config struct {
	PollInterval string            `json:"pollInterval,omitempty"`
	Instances    []config.Instance `json:"instances,omitempty"`
}

type Provider struct {
	config       *Config
	aggregators  []aggregator.Aggregator
	pollInterval time.Duration

	cancel func()
}

func CreateConfig() *Config {
	return &Config{}
}

func New(ctx context.Context, config *Config, name string) (*Provider, error) {
	log.Info("New Traefik Aggregator plugin with name %s and config %+v", name, config)
	pollInterval, err := time.ParseDuration(config.PollInterval)

	if err != nil {
		return nil, fmt.Errorf("error while parsing poll interval: %w", err)
	}

	provider := Provider{
		config:       config,
		pollInterval: pollInterval,
	}

	return &provider, nil
}

func (p *Provider) Init() error {
	log.Info("Initializing Traefik Aggregator plugin with config %+v", p.config)

	p.aggregators = make([]aggregator.Aggregator, len(p.config.Instances))

	for i, instance := range p.config.Instances {
		log.Info("Instance %d: %+v", i, instance)
		fetcher := aggregator.NewFetcher(instance)
		converter := aggregator.NewConverter(instance)

		p.aggregators[i] = aggregator.NewAggregator(fetcher, converter, instance)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("Recovered from panic: %v", r)
			}
		}()

		for true {
			time.Sleep(2 * time.Second)
			for _, aggregator := range p.aggregators {
				aggregator.Run()
			}
		}
	}()

	return nil
}

func (p *Provider) Provide(cfgChan chan<- json.Marshaler) error {
	ctx, cancel := context.WithCancel(context.Background())

	p.cancel = cancel

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Recovered from panic: %v", err)
			}
		}()

		p.run(ctx, cfgChan)
	}()
	return nil
}

func (p *Provider) run(ctx context.Context, cfgChan chan<- json.Marshaler) {
	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			configs := []*dynamic.Configuration{}

			for _, aggregator := range p.aggregators {
				config, err := aggregator.Run()

				if err != nil {
					log.Error("Error while running aggregator: %v", err)
					continue
				}

				configs = append(configs, config)
			}

			configuration := aggregator.JoinConfigs(configs)

			cfgChan <- &dynamic.JSONPayload{Configuration: configuration}

		case <-ctx.Done():
			return
		}
	}
}

func (p *Provider) Stop() error {
	return nil
}
