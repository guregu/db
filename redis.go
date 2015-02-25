package db

import (
	"golang.org/x/net/context"
	"gopkg.in/redis.v2"
)

type rediskey string

// Redis retrieves the *redis.Client with the given name or nil.
func Redis(ctx context.Context, name string) *redis.Client {
	client, _ := ctx.Value(rediskey(name)).(*redis.Client)
	return client
}

// WithRedis returns a new context containing the given *redis.Client
func WithRedis(ctx context.Context, name string, client *redis.Client) context.Context {
	key := rediskey(name)
	if idx := redisIndexFrom(ctx); idx != nil {
		idx[key] = client
	} else {
		idx = map[rediskey]*redis.Client{key: client}
		ctx = withRedisIndex(ctx, idx)
	}
	return context.WithValue(ctx, rediskey(name), client)
}

// OpenRedis opens a Redis connection and returns a new context or panics.
func OpenRedis(ctx context.Context, name string, options *redis.Options) context.Context {
	client := redis.NewClient(options)
	return WithRedis(ctx, name, client)
}

// OpenFailoverRedis opens a failover Redis connection and returns a new context or panics.
func OpenFailoverRedis(ctx context.Context, name string, options *redis.FailoverOptions) context.Context {
	client := redis.NewFailoverClient(options)
	return WithRedis(ctx, name, client)
}

// CloseRedis closes the specified Redis connection, panciking if Close returns an error.
// CloseRedis will do nothing if the given Redis connection does not exist.
func CloseRedis(ctx context.Context, name string) context.Context {
	client := Redis(ctx, name)
	if client == nil {
		return ctx
	}

	if err := client.Close(); err != nil {
		panic(err)
	}
	return removeRedis(ctx, name)
}

// CloseRedisAll closes all open Redis connections and returns a new context without them.
func CloseRedisAll(ctx context.Context) context.Context {
	if idx := redisIndexFrom(ctx); idx != nil {
		for name, _ := range idx {
			ctx = CloseRedis(ctx, string(name))
		}
	}
	return ctx
}

func removeRedis(ctx context.Context, name string) context.Context {
	key := rediskey(name)
	if idx := redisIndexFrom(ctx); idx != nil {
		delete(idx, key)
	}
	return context.WithValue(ctx, key, nil)
}

func redisIndexFrom(ctx context.Context) map[rediskey]*redis.Client {
	idx, _ := ctx.Value(redisIndex).(map[rediskey]*redis.Client)
	return idx
}

func withRedisIndex(ctx context.Context, idx map[rediskey]*redis.Client) context.Context {
	return context.WithValue(ctx, redisIndex, idx)
}
