package application

import (
	"strings"

	"github.com/MadddinTribleD/traefikaggregator/pkg/config"
	"github.com/MadddinTribleD/traefikaggregator/pkg/models"
	"github.com/traefik/genconf/dynamic"
)

type RouterConfigBuilder struct {
	serviceName         string
	certResolverMapping map[string]string
	config              config.RouterConfig
	allowedEndpoints    []string
}

func NewRouterConfigBuilder(
	serviceName string,
	certResolverMapping map[string]string,
	config config.RouterConfig,
	allowedEndpoints []string,
) *RouterConfigBuilder {
	return &RouterConfigBuilder{
		serviceName:         serviceName,
		certResolverMapping: certResolverMapping,
		config:              config,
		allowedEndpoints:    allowedEndpoints,
	}
}

func (c *RouterConfigBuilder) BuildHttpConfig(routers []models.Router) map[string]*dynamic.Router {
	routerMap := map[string]*dynamic.Router{}

	for _, router := range c.filterRouters(routers) {
		routerName := c.buildRouterName(router)
		routerMap[routerName] = &dynamic.Router{
			EntryPoints: c.config.EntryPoints,
			Middlewares: c.config.Middlewares,
			Priority:    c.config.Priority,
			Rule:        router.Rule,
			Service:     c.serviceName,
		}

		if c.config.TLS != nil {
			if certResolver, ok := c.certResolverMapping[router.TLS.CertResolver]; ok {
				resolver := certResolver
				routerMap[routerName].TLS = &dynamic.RouterTLSConfig{
					Options:      c.config.TLS.Options,
					Domains:      c.config.TLS.Domains,
					CertResolver: resolver,
				}
			}
		}
	}

	return routerMap
}

func (c *RouterConfigBuilder) filterRouters(routers []models.Router) []models.Router {
	filteredRouters := []models.Router{}

	for _, router := range routers {
		for _, allowedEndpoint := range c.allowedEndpoints {
			found := false
			for _, endpoint := range router.EntryPoints {
				if endpoint == allowedEndpoint {
					filteredRouters = append(filteredRouters, router)
					found = true
					break
				}
			}

			if found {
				break
			}
		}
	}

	return filteredRouters
}

func (c *RouterConfigBuilder) buildRouterName(router models.Router) string {
	return strings.Replace(router.Name, "@docker", "", -1)
}

// func (c *RouterConfigBuilder) buildRouters(routers []models.Router) map[string]*dynamic.Router {
// 	routerMap := map[string]*dynamic.Router{}

// 	for _, router := range routers {
// 		routerMap[c.buildRouterName(router)] = &router.Router
// 	}

// 	return routerMap
// }
