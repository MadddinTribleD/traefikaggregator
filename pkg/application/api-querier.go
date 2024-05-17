package application

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MadddinTribleD/traefikaggregator/pkg/models"
	"github.com/traefik/genconf/dynamic"
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

	routers, err := a.unmarshalBody(body)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling routers: %w", err)
	}

	return routers, nil
}

func (a *ApiQuerier) unmarshalBody(body []byte) ([]models.Router, error) {
	dynamicRouter := []dynamic.Router{}
	if err := json.Unmarshal(body, &dynamicRouter); err != nil {
		return nil, fmt.Errorf("error while unmarshalling routers: %w", err)
	}

	names := []struct {
		Name string `json:"name"`
	}{}
	if err := json.Unmarshal(body, &names); err != nil {
		return nil, fmt.Errorf("error while unmarshalling router names: %w", err)
	}

	routers := []models.Router{}

	for i, router := range dynamicRouter {
		routers = append(routers, models.Router{
			Router: router,
			Name:   names[i].Name,
		})
	}

	return routers, nil

}
