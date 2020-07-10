package mid

import (
	"context"
	"expvar"
	"net/http"
	"runtime"

	"github.com/sreejeet/garagesale/internal/platform/web"
)

// m is a global structure containing all the applcation
// metric varaiables that need to be minitored.
var m = struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	errors     *expvar.Int
}{
	goroutines: expvar.NewInt("goroutines"),
	requests:   expvar.NewInt("requests"),
	errors:     expvar.NewInt("errors"),
}

// Metrics middleware updates the application metrics
// each time any request is processesd.
func Metrics() web.Middleware {

	f := func(before web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// Exeute the handler before this middleware here.
			err := before(ctx, w, r)

			// Increment the request counter.
			m.requests.Add(1)

			// Update the count for the number of active goroutines every 100 requests.
			if m.requests.Value()%100 == 0 {
				m.goroutines.Set(int64(runtime.NumGoroutine()))
			}

			// Increment the errors counter if an error occurred on this request.
			if err != nil {
				m.errors.Add(1)
			}

			// Return the error so it can be handled further up the chain.
			return err
		}

		return h
	}

	return f
}
