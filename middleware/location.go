package middleware

import (
	"fmt"
	"time"

	"github.com/elfincafe/bismarck"
)

func Location(loc *time.Location) bismarck.MiddlewareFunc {
	return func(next bismarck.HandlerFunc) bismarck.HandlerFunc {
		return func(ctx *bismarck.Context) error {
			fmt.Println("Location Middleware Executing")
			ctx.SetLocation(loc)
			return next(ctx)
		}
	}
}
