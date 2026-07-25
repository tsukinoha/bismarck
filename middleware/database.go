package middleware

import (
	"database/sql"

	"github.com/tsukinoha/bismarck"
)

type (
	DbConfig struct {
	}
)

func Database(db *sql.DB, config *DbConfig) bismarck.MiddlewareFunc {
	return func(next bismarck.HandlerFunc) bismarck.HandlerFunc {
		return func(ctx *bismarck.Context) error {
			ctx.SetDatabase(db)
			return next(ctx)
		}
	}
}
