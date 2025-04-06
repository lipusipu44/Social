package cache

import (
	"context"
	"github.com/lipusipu44/Social/internal/store"
	"github.com/redis/go-redis/v9"
)

type CacheStorage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}
}

func NewRedisStorage(rc *redis.Client) CacheStorage {
	return CacheStorage{
		Users: &UserClient{
			rc: rc,
		},
	}
}
