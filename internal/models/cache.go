package models

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) 		error
	SetNX(ctx context.Context, key string, value any, ttl time.Duration)	*redis.BoolCmd
	Get(ctx context.Context, key string)									(any, error)
	Del(ctx context.Context, key string)									error
}