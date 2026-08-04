package middleware

import (
	"github.com/tsukinoha/bismarck"
)

func Warehouse(key string, value any) bismarck.MiddlewareFunc {
	return func(next bismarck.HandlerFunc) bismarck.HandlerFunc {
		return func(ctx *bismarck.Context) error {
			ctx.Set(key, value)
			return next(ctx)
		}
	}
}
