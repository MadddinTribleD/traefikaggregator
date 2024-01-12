package config

type Instance struct {
	ServiceName                 string                        `json:"serviceName"`
	Urls                        []string                      `json:"urls"`
	ApiEndpoint                 string                        `json:"apiEndpoint"`
	CertificatesResolverMapping []CertificatesResolverMapping `json:"certificatesResolverMapping"`
	AllowedEndpoints            []string                      `json:"allowedEndpoints"`
	EntryPoint                  string                        `json:"entryPoint"`
}

type CertificatesResolverMapping struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}
