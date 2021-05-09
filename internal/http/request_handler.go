package http

import "net/http"

type Middleware func(handler http.Handler) http.Handler

type Route struct {
	Pattern string
	Handler http.Handler
}

func RequestHandler(routes []Route, middlewares ...Middleware) http.Handler {
	mux := http.NewServeMux()

	for _, route := range routes {
		mux.Handle(route.Pattern, route.Handler)
	}

	handler := http.Handler(mux)

	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
