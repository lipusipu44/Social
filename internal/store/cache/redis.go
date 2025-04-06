package cache

import (
	"github.com/redis/go-redis/v9"
)

//NewRedislient
/*
Initialization of Redis like we did for DB in storage.go
*/
func NewRedislient(addr, pw string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})
}
