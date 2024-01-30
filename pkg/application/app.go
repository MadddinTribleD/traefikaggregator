package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/traefik/genconf/dynamic"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/config"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/models"
)

type App interface {
	Run(ctx context.Context, configChannel chan<- json.Marshaler) error
}

type routerConfigBuilderWithApi struct {
	routerConfigBuilder *RouterConfigBuilder
	apiEndpoint         string
}

type app struct {
	pollInterval        time.Duration
	apiQueriers         map[string]*ApiQuerier
	routerConfigBuilder map[string]*routerConfigBuilderWithApi
	serviceConfig       map[string]*dynamic.Service
}

func NewApp(config *config.Config) (App, error) {
	apiQueriers := map[string]*ApiQuerier{}
	routerConfigBuilders := map[string]*routerConfigBuilderWithApi{}

	pollInterval, err := time.ParseDuration(config.PollInterval)

	if err != nil {
		return &app{}, fmt.Errorf("error while parsing poll interval: %w", err)
	}

	serviceConfig := NewServiceConfigBuilder(config.Instances).BuildServiceConfig()

	for _, instance := range config.Instances {
		if _, ok := apiQueriers[instance.ApiEndpoint]; !ok {
			apiQueriers[instance.ApiEndpoint] = NewApiQuerier(instance.ApiEndpoint)
		}

		routerConfigBuilder := NewRouterConfigBuilder(
			instance.Service.Name,
			instance.CertResolverMapping,
			instance.Router,
			instance.AllowedEndpoints,
		)

		routerConfigBuilders[instance.Service.Name] = &routerConfigBuilderWithApi{
			routerConfigBuilder: routerConfigBuilder,
			apiEndpoint:         instance.ApiEndpoint,
		}
	}

	return &app{
		pollInterval:        pollInterval,
		apiQueriers:         apiQueriers,
		routerConfigBuilder: routerConfigBuilders,
		serviceConfig:       serviceConfig,
	}, nil
}

func (a *app) Run(ctx context.Context, configChannel chan<- json.Marshaler) error {
	ticker := time.NewTicker(a.pollInterval)
	defer ticker.Stop()

	var oldConfigByte []byte

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			config, err := a.run(ctx)

			if err != nil {
				return err
			}

			configByte, _ := config.MarshalJSON()

			if bytes.Equal(configByte, oldConfigByte) {
				continue
			}

			oldConfigByte = configByte

			configChannel <- config
		}
	}
}

func (a *app) run(ctx context.Context) (json.Marshaler, error) {
	apiRouters := map[string][]models.Router{}

	for apiEndpoint, apiQuerier := range a.apiQueriers {
		routers, err := apiQuerier.QueryHttpRouter(ctx)

		if err != nil {
			return nil, err
		}

		apiRouters[apiEndpoint] = routers
	}

	allRouterConfig := map[string]*dynamic.Router{}

	for _, routerConfigBuilder := range a.routerConfigBuilder {
		apiRouter := apiRouters[routerConfigBuilder.apiEndpoint]
		routerConfig := routerConfigBuilder.routerConfigBuilder.BuildHttpConfig(apiRouter)

		for routerName, router := range routerConfig {
			allRouterConfig[routerName] = router
		}
	}

	return &dynamic.JSONPayload{
		Configuration: &dynamic.Configuration{
			HTTP: &dynamic.HTTPConfiguration{
				Routers:  allRouterConfig,
				Services: a.serviceConfig,
			},
		},
	}, nil
}
