package mid

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/sreejeet/garagesale/internal/platform/web"
	"go.opencensus.io/trace"
)

// Logger middleware logs info for each requets in the format
// (200) GET /foo -> IP ADDR (latency)
func Logger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			ctx, span := trace.StartSpan(ctx, "internal.mid.RequestLogger")
			defer span.End()

			// Get keys from the context, return error in case the value is not found.
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return errors.New("web value missing from context")
			}

			err := before(ctx, w, r)

			log.Printf("(%d) : %s %s -> %s (%s)",
				v.StatusCode,
				r.Method, r.URL.Path,
				r.RemoteAddr, time.Since(v.Start),
			)

			// Return the error so it can be handled further up the chain.
			return err
		}

		return h
	}

	return f
}
