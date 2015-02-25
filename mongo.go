package db

import (
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

type mongokey string

// Mongo retrieves the *mgo.Session with the given name or nil.
func MongoDB(ctx context.Context, name string) *mgo.Session {
	db, _ := ctx.Value(mongokey(name)).(*mgo.Session)
	return db
}

// WithMongoDB returns a new context containing the given *mgo.Session
func WithMongoDB(ctx context.Context, name string, db *mgo.Session) context.Context {
	key := mongokey(name)
	if idx := mongoIndexFrom(ctx); idx != nil {
		idx[key] = db
	} else {
		idx = map[mongokey]*mgo.Session{key: db}
		ctx = withMongoIndex(ctx, idx)
	}
	return context.WithValue(ctx, mongokey(name), db)
}

// OpenMongoDB opens a Mongo connection and returns a new context or panics.
func OpenMongoDB(ctx context.Context, name, url string) context.Context {
	db, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	return WithMongoDB(ctx, name, db)
}

// CloseMongoDB closes the specified Mongo connection, panciking if Close returns an error.
// CloseMongoDB will do nothing if the given Mongo connection does not exist.
func CloseMongoDB(ctx context.Context, name string) context.Context {
	db := MongoDB(ctx, name)
	if db == nil {
		return ctx
	}

	db.Close()
	return removeMongo(ctx, name)
}

// CloseMongoDBAll closes all open Mongo connections and returns a new context without them.
func CloseMongoDBAll(ctx context.Context) context.Context {
	if idx := mongoIndexFrom(ctx); idx != nil {
		for name, _ := range idx {
			ctx = CloseMongoDB(ctx, string(name))
		}
	}
	return ctx
}

func removeMongo(ctx context.Context, name string) context.Context {
	key := mongokey(name)
	if idx := mongoIndexFrom(ctx); idx != nil {
		delete(idx, key)
	}
	return context.WithValue(ctx, key, nil)
}

func mongoIndexFrom(ctx context.Context) map[mongokey]*mgo.Session {
	idx, _ := ctx.Value(mongoIndex).(map[mongokey]*mgo.Session)
	return idx
}

func withMongoIndex(ctx context.Context, idx map[mongokey]*mgo.Session) context.Context {
	return context.WithValue(ctx, mongoIndex, idx)
}
