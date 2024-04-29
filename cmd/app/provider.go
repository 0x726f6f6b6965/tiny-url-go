package main

import (
	"context"

	"github.com/0x726f6f6b6965/tiny-url-go/cmd/api"
	"github.com/0x726f6f6b6965/tiny-url-go/internal/config"
	"github.com/0x726f6f6b6965/tiny-url-go/internal/service"
	"github.com/0x726f6f6b6965/tiny-url-go/internal/storage"
	"github.com/0x726f6f6b6965/tiny-url-go/utils"
	"github.com/google/wire"
)

var applicationSet = wire.NewSet(dynamoDBSet, loggerSet, sequencerSet, service.NewTinyURLService, api.NewShortenAPI)

var sequencerSet = wire.NewSet(sequencerConfig, initSequencer)

var loggerSet = wire.NewSet(logConfig, utils.NewLogger)

func sequencerConfig(cfg *config.AppConfig) *config.SequencerConfig {
	return &cfg.Sequencer
}

func initSequencer(cfg *config.SequencerConfig) (utils.Sequencer, error) {
	return utils.NewSequencer(cfg.NodeID, cfg.Start)
}

func dynamoDBSet(ctx context.Context, cfg *config.AppConfig) (utils.Storage, error) {
	if cfg.Env == "dev" {
		return storage.NewDevDynamoDB(ctx, &cfg.Storage)
	} else {
		return storage.NewDynamoDB(ctx, &cfg.Storage)
	}
}

func logConfig(cfg *config.AppConfig) *config.LogConfig {
	return &cfg.Log
}
