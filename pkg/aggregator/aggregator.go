package aggregator

import (
	"github.com/traefik/genconf/dynamic"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/config"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/models"
)

type Aggregator interface {
	Run() (*dynamic.Configuration, error)
}

func NewAggregator(fetcher Fetcher, converter Converter, config config.Instance) Aggregator {
	allowedEndpoints := map[string]interface{}{}

	for _, endpoint := range config.AllowedEndpoints {
		allowedEndpoints[endpoint] = nil
	}

	return &aggregator{
		fetcher,
		converter,
		allowedEndpoints,
		config.ServiceName,
		config.Urls,
	}
}

type aggregator struct {
	fetcher   Fetcher
	converter Converter

	allowedEndpoints map[string]interface{}
	serviceName      string
	urls             []string
}

func (a *aggregator) Run() (*dynamic.Configuration, error) {
	routers, err := a.fetcher.FetchRouters()

	if err != nil {
		return nil, err
	}

	routers = a.filter(routers)

	convertedRouters := []models.Router{}
	for _, router := range routers {

		if convertedRouter, err := a.converter.Convert(router); err != nil {
			return nil, err
		} else {
			convertedRouters = append(convertedRouters, *convertedRouter)
		}
	}

	config := &dynamic.Configuration{
		HTTP: &dynamic.HTTPConfiguration{
			Routers:  makeRouterMap(convertedRouters),
			Services: a.makeServiceMap(),
		},
	}

	return config, nil
}

func (a *aggregator) makeServiceMap() map[string]*dynamic.Service {
	serviceMap := map[string]*dynamic.Service{}

	serviceMap[a.serviceName] = &dynamic.Service{
		LoadBalancer: &dynamic.ServersLoadBalancer{
			Servers: makeServerSlice(a.urls),
		},
	}

	return serviceMap
}

func makeServerSlice(urls []string) []dynamic.Server {
	servers := []dynamic.Server{}

	for _, url := range urls {
		servers = append(servers, dynamic.Server{
			URL: url,
		})
	}

	return servers
}

func makeRouterMap(routers []models.Router) map[string]*dynamic.Router {
	routerMap := map[string]*dynamic.Router{}

	for _, router := range routers {
		routerMap[router.Name] = &router.Router
	}

	return routerMap
}

func (a *aggregator) filter(routers []models.Router) []models.Router {
	filteredRouters := []models.Router{}

	for _, router := range routers {
		for _, endpoint := range router.EntryPoints {
			if _, ok := a.allowedEndpoints[endpoint]; ok {
				filteredRouters = append(filteredRouters, router)
				break
			}
		}
	}

	return filteredRouters
}
