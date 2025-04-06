package cache

import (
	"context"
	"github.com/lipusipu44/Social/internal/store"
	"github.com/redis/go-redis/v9"
)

type UserClient struct {
	rc *redis.Client
}

func (u *UserClient) Get(ctx context.Context, db int64) (*store.User, error) {
	return nil, nil
}

func (u *UserClient) Set(ctx context.Context, usr *store.User) error {
	return nil
}
