package http

import "context"

type (
	Handler    func(context.Context, Request) Response
	Middleware func(Handler) Handler
)
