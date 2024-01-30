package main

import (
	"context"
	"encoding/json"

	"version.gafert.org/MadddinTribleD/traefikaggregator"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/config"
	"version.gafert.org/MadddinTribleD/traefikaggregator/pkg/log"
)

func main() {
	ctx := context.Background()

	config := &config.Config{
		PollInterval: "2s",
		Instances:    []config.TraefikInstance{},
	}
	plugin, _ := traefikaggregator.New(ctx, config, "traefikaggregator")

	plugin.Init()

	cfgChan := make(chan json.Marshaler)

	err := plugin.Provide(cfgChan)
	if err != nil {
		panic(err)
	}

	data := <-cfgChan

	log.Info("Configuration: %s", data)
}
