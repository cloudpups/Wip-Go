package main

import (
	"fmt"
	"time"

	"net/http"

	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
)

func runApp(config Config, logger zerolog.Logger) {
	cc, err := githubapp.NewDefaultCachingClientCreator(
		config.Github,
		githubapp.WithClientUserAgent("WIP-go"),
		githubapp.WithClientTimeout(3*time.Second),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
	)

	if err != nil {
		panic(err)
	}

	prStatusHandler := &PRStatusHandler{
		ClientCreator: cc,
	}

	webhookHandler := githubapp.NewDefaultEventDispatcher(config.Github, prStatusHandler)

	http.Handle(githubapp.DefaultWebhookRoute, webhookHandler)

	addr := fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port)

	logger.Info().Msgf("Starting server on %s...", addr)

	err = http.ListenAndServe(addr, nil)

	if err != nil {
		panic(err)
	}
}
