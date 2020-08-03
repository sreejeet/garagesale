package mid

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/sreejeet/garagesale/internal/platform/web"
	"go.opencensus.io/trace"
)

// Panics is used to recover and convert panics to errors so that
// they can be reported in the metrics and handled in the Errors.
func Panics(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(after web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {

			ctx, span := trace.StartSpan(ctx, "internal.mid.Panics")
			defer span.End()

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}

			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if r := recover(); r != nil {
					err = errors.Errorf("panic: %v", r)

					// Log the Go stack trace for this panic'd goroutine.
					log.Printf("%s :\n%s", v.TraceID, debug.Stack())
				}
			}()

			// Call the next Handler and set its return value in the err variable.
			return after(ctx, w, r)
		}

		return h
	}

	return f
}
