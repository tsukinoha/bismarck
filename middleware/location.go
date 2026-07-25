package middleware

import (
	"time"

	"github.com/tsukinoha/bismarck"
)

func Location(loc *time.Location) bismarck.MiddlewareFunc {
	return func(next bismarck.HandlerFunc) bismarck.HandlerFunc {
		return func(ctx *bismarck.Context) error {
			ctx.SetLocation(loc)
			return next(ctx)
		}
	}
}
