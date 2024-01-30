package application

import (
	"github.com/traefik/genconf/dynamic"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/config"
)

type ServiceConfigBuilder struct {
	instances []config.TraefikInstance
}

func NewServiceConfigBuilder(instances []config.TraefikInstance) *ServiceConfigBuilder {
	return &ServiceConfigBuilder{
		instances: instances,
	}
}

func (c *ServiceConfigBuilder) BuildServiceConfig() map[string]*dynamic.Service {
	serviceMap := map[string]*dynamic.Service{}

	for _, instance := range c.instances {
		inst := instance
		service := dynamic.Service{
			LoadBalancer: &inst.Service.LoadBalancer,
		}
		serviceMap[inst.Service.Name] = &service
	}

	return serviceMap
}
