package http

import "context"

// Implementing this idea https://hackernoon.com/simple-http-middleware-with-go-79a4ad62889b.
// The mw is an array consist of Middleware type.
func ChainMiddleware(mw ...Middleware) Middleware {
	return func(final Handler) Handler {
		return func(ctx context.Context, req Request) Response {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}

			// last middleware
			return last(ctx, req)
		}
	}
}
