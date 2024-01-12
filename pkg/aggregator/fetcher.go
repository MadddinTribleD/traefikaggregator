package aggregator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/config"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/models"
)

type Fetcher interface {
	FetchRouters() ([]models.Router, error)
}

func NewFetcher(instanceConfig config.Instance) Fetcher {
	return &fetcher{instanceConfig}
}

type fetcher struct {
	config config.Instance
}

func (a *fetcher) FetchRouters() ([]models.Router, error) {
	res, err := http.Get(fmt.Sprintf("%s/http/routers", a.config.ApiEndpoint))

	if err != nil {
		return nil, fmt.Errorf("error while fetching routers: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while fetching routers: %s", res.Status)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("error while reading response body: %w", err)
	}

	routers := []models.Router{}
	if err := json.Unmarshal(body, &routers); err != nil {
		return nil, fmt.Errorf("error while unmarshalling routers: %w", err)
	}

	return routers, nil
}
