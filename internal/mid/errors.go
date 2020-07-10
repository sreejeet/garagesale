package mid

import (
	"context"
	"log"
	"net/http"

	"github.com/sreejeet/garagesale/internal/platform/web"
)

// Errors handles errors from the middleware chains and
// respond to normal application errors in a uniform manner.
// Any unexpected errors will be responded with a 5xx status code
// and will be logged.
func Errors(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// Run the handler chain and catch any propagated error.
			if err := before(ctx, w, r); err != nil {

				// Log the error.
				log.Printf("ERROR : %v", err)

				// Respond to the error.
				if err := web.RespondError(ctx, w, err); err != nil {
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
