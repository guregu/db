package db

import "golang.org/x/net/context"

type indexkey int

const (
	sqlIndex indexkey = iota + 1
	redisIndex
	mongoIndex
)

// Close closes all connections of all kinds and returns a new context without them.
func Close(ctx context.Context) context.Context {
	ctx = CloseSQLAll(ctx)
	ctx = CloseRedisAll(ctx)
	return ctx
}
