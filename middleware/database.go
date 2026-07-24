package middleware

import (
	"database/sql"
	"fmt"

	"github.com/elfincafe/bismarck"
)

type (
	DbConfig struct {
	}
)

func Database(db *sql.DB, config *DbConfig) bismarck.MiddlewareFunc {
	return func(next bismarck.HandlerFunc) bismarck.HandlerFunc {
		return func(ctx *bismarck.Context) error {
			fmt.Println("Database Middleware Executing")
			ctx.SetDatabase(db)
			return next(ctx)
		}
	}
}
