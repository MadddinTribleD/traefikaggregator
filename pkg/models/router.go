package models

import "github.com/traefik/genconf/dynamic"

type Router struct {
	dynamic.Router

	Name string `json:"name,omitempty"`
}
