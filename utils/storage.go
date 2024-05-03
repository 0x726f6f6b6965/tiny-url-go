package utils

import "context"

type Storage interface {
	Save(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
	Update(ctx context.Context, key string, value interface{}, updateMask []string) error
}
