package mid

import (
	"context"
	"log"
	"net/http"

	"github.com/sreejeet/garagesale/internal/platform/web"
	"go.opencensus.io/trace"
)

// Errors handles errors from the middleware chains and
// respond to normal application errors in a uniform manner.
// Any unexpected errors will be responded with a 5xx status code
// and will be logged.
func Errors(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			ctx, span := trace.StartSpan(ctx, "internal.mid.Errors")
			defer span.End()

			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web value missing from context")
			}
			// Run the handler chain and catch any propagated error.
			if err := before(ctx, w, r); err != nil {

				// Log the error.
				log.Printf("%s : ERROR : %+v", v.TraceID, err)

				// Respond to the error.
				if err := web.RespondError(ctx, w, err); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shutdown the service.
				if ok := web.IsShutdown(err); ok {
					return err
				}

			}

			// Return nil to indicate the error has been handled.
			return nil
		}

		return h
	}

	return f
}
