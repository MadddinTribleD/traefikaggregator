package config

import (
	"github.com/traefik/genconf/dynamic"
	"github.com/traefik/genconf/dynamic/types"
)

type Config struct {
	PollInterval string            `json:"pollInterval,omitempty"`
	Instances    []TraefikInstance `json:"instances,omitempty"`
}

type TraefikInstance struct {
	ApiEndpoint         string            `json:"apiEndpoint"`
	AllowedEndpoints    []string          `json:"allowedEndpoints"`
	Router              RouterConfig      `json:"router,omitempty"`
	CertResolverMapping map[string]string `json:"certResolverMapping"`
	Service             ServiceConfig     `json:"service,omitempty"`
}

type ServiceConfig struct {
	Name         string                      `json:"name"`
	LoadBalancer dynamic.ServersLoadBalancer `json:"loadBalancer,omitempty"`
}

type RouterConfig struct {
	EntryPoints []string `json:"entryPoints,omitempty"`
	Middlewares []string `json:"middlewares,omitempty"`
	// omit service and rule because we build them by our self
	Priority int              `json:"priority,omitempty"`
	TLS      *RouterTLSConfig `json:"tls,omitempty"`
}

type RouterTLSConfig struct {
	Enabled bool           `json:"enabled,omitempty"`
	Options string         `json:"options,omitempty"`
	Domains []types.Domain `json:"domains,omitempty"`
}
