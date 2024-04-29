//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/0x726f6f6b6965/tiny-url-go/cmd/api"
	"github.com/0x726f6f6b6965/tiny-url-go/internal/config"
	"github.com/google/wire"
)

func initApplication(ctx context.Context, cfg *config.AppConfig) (application *api.ShortenAPI, err error) {
	panic(wire.Build(applicationSet))
}
