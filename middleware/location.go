package middleware

import (
	"time"

	"github.com/elfincafe/bismarck"
)

func Location(loc *time.Location) bismarck.MiddlewareFunc {
	return func(next bismarck.HandlerFunc) bismarck.HandlerFunc {
		return func(ctx *bismarck.Context) error {
			ctx.loc = loc
			return next(ctx)
		}
	}
}
