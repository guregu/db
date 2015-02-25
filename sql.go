package db

import (
	"database/sql"

	"golang.org/x/net/context"
)

type sqlkey string

// SQL retrieves the *sql.DB with the given name or nil.
func SQL(ctx context.Context, name string) *sql.DB {
	db, _ := ctx.Value(sqlkey(name)).(*sql.DB)
	return db
}

// WithSQL returns a new context containing the given *sql.DB
func WithSQL(ctx context.Context, name string, db *sql.DB) context.Context {
	key := sqlkey(name)
	if idx := sqlIndexFrom(ctx); idx != nil {
		idx[key] = db
	} else {
		idx = map[sqlkey]*sql.DB{key: db}
		ctx = withSQLIndex(ctx, idx)
	}
	return context.WithValue(ctx, sqlkey(name), db)
}

// OpenSQL opens a SQL connection and returns a new context or panics.
func OpenSQL(ctx context.Context, name, driver, dataSource string) context.Context {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		panic(err)
	}
	return WithSQL(ctx, name, db)
}

// CloseSQL closes the specified SQL connection, panciking if Close returns an error.
// CloseSQL will do nothing if the given SQL connection does not exist.
func CloseSQL(ctx context.Context, name string) context.Context {
	db := SQL(ctx, name)
	if db == nil {
		return ctx
	}

	if err := db.Close(); err != nil {
		panic(err)
	}
	return removeSQL(ctx, name)
}

// CloseSQLAll closes all open SQL connections and returns a new context without them.
func CloseSQLAll(ctx context.Context) context.Context {
	if idx := sqlIndexFrom(ctx); idx != nil {
		for name, _ := range idx {
			ctx = CloseSQL(ctx, string(name))
		}
	}
	return ctx
}

func removeSQL(ctx context.Context, name string) context.Context {
	key := sqlkey(name)
	if idx := sqlIndexFrom(ctx); idx != nil {
		delete(idx, key)
	}
	return context.WithValue(ctx, key, nil)
}

func sqlIndexFrom(ctx context.Context) map[sqlkey]*sql.DB {
	idx, _ := ctx.Value(sqlIndex).(map[sqlkey]*sql.DB)
	return idx
}

func withSQLIndex(ctx context.Context, idx map[sqlkey]*sql.DB) context.Context {
	return context.WithValue(ctx, sqlIndex, idx)
}
