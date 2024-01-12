package aggregator

import (
	"fmt"
	"strings"

	"github.com/traefik/genconf/dynamic"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/config"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/models"
)

type Converter interface {
	Convert(router models.Router) (*models.Router, error)
}

func NewConverter(config config.Instance) Converter {
	certificatesResolverMapping := map[string]string{}

	for _, certificateResolver := range config.CertificatesResolverMapping {
		certificatesResolverMapping[certificateResolver.Source] = certificateResolver.Destination
	}

	return &converter{
		config:                      config,
		certificatesResolverMapping: certificatesResolverMapping,
	}
}

type converter struct {
	config                      config.Instance
	certificatesResolverMapping map[string]string
}

func (c *converter) Convert(router models.Router) (*models.Router, error) {

	convertedRouter := models.Router{
		Router: dynamic.Router{
			EntryPoints: router.EntryPoints,
			Rule:        router.Rule,
			Service:     c.config.ServiceName,
		},

		Name: strings.Replace(router.Name, "@docker", "", -1),
	}

	if router.TLS != nil {
		if certResolver, ok := c.certificatesResolverMapping[router.TLS.CertResolver]; ok {
			convertedRouter.TLS = &dynamic.RouterTLSConfig{
				CertResolver: certResolver,
			}
		} else {
			return nil, fmt.Errorf("could not find certificate resolver %s", router.TLS.CertResolver)
		}
	}

	return &convertedRouter, nil
}
