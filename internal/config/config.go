package config

import "time"

type AppConfig struct {
	Env       string          `yaml:"env" mapstructure:"env" cobra-usage:"the application environment" cobra-default:"dev"`
	Host      string          `yaml:"host" mapstructure:"host" validate:"required" cobra-usage:"the application host" cobra-default:"localhost"`
	Port      uint64          `yaml:"port" mapstructure:"port" validate:"required,gte=0" cobra-usage:"the application port" cobra-default:"8080"`
	Log       LogConfig       `yaml:"log" mapstructure:"log"`
	TableName string          `yaml:"table-name" mapstructure:"table-name" cobra-usage:"the dynamodb table name" cobra-default:""`
	Expire    time.Duration   `yaml:"expire" mapstructure:"expire"`
	Sequencer SequencerConfig `yaml:"sequencer" mapstructure:"sequencer"`
	Storage   StorageConfig   `yaml:"storage" mapstructure:"storage"`
}

type SequencerConfig struct {
	NodeID int64     `yaml:"node-id" mapstructure:"node-id" validate:"omitempty,gte=0" cobra-usage:"the node id" cobra-default:"1"`
	Start  time.Time `yaml:"start" mapstructure:"start" validate:"required" cobra-usage:"the start time" cobra-default:""`
}

type StorageConfig struct {
	Region string `yaml:"region" mapstructure:"region" validate:"required" cobra-usage:"the storage region" cobra-default:"us-east-1"`
	Host   string `yaml:"host" mapstructure:"host" validate:"omitempty" cobra-usage:"the storage host" cobra-default:"localhost"`
	Port   uint64 `yaml:"port" mapstructure:"port" validate:"omitempty,gte=0" cobra-usage:"the storage port" cobra-default:"8686"`
}

type LogConfig struct {
	Level            int    `yaml:"level" mapstructure:"level" validate:"omitempty,gte=-1,lte=5" cobra-usage:"the application log level" cobra-default:"1"`
	TimeFormat       string `yaml:"time-format" mapstructure:"time-format" cobra-usage:"the application log time format" cobra-default:"2006-01-02T15:04:05Z07:00"`
	TimestampEnabled bool   `yaml:"timestamp-enabled" mapstructure:"timestamp-enabled" cobra-usage:"specify if the timestamp is enabled"  cobra-default:"false"`
	ServiceName      string `yaml:"service-name" mapstructure:"service-name" cobra-usage:"the application service name" cobra-default:""`
}
