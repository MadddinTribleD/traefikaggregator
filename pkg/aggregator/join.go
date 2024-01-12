package aggregator

import (
	"github.com/traefik/genconf/dynamic"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/log"
)

func JoinConfigs(configs []*dynamic.Configuration) *dynamic.Configuration {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Recovered from panic: %v", err)
		}
	}()

	joined := dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers:  map[string]*dynamic.Router{},
			Services: map[string]*dynamic.Service{},
		},
	}

	for _, config := range configs {
		for routerName, router := range config.HTTP.Routers {
			joined.HTTP.Routers[routerName] = router
		}

		for serviceName, service := range config.HTTP.Services {
			joined.HTTP.Services[serviceName] = service
		}
	}

	return &joined
}
