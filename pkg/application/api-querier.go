package application

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/models"
)

type ApiQuerier struct {
	apiEndpoint string
}

func NewApiQuerier(apiEndpoint string) *ApiQuerier {
	return &ApiQuerier{
		apiEndpoint: apiEndpoint,
	}
}

func (a *ApiQuerier) QueryHttpRouter(ctx context.Context) ([]models.Router, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/http/routers", a.apiEndpoint), nil)
	client := &http.Client{}
	res, err := client.Do(req)

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
