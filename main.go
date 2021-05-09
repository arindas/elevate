package main

import (
	"flag"

	"github.com/arindas/elevate/internal/app"
	"github.com/arindas/elevate/internal/http"
)

var baseDirectory = flag.String("base_dir", ".", "Base directory for storing files.")

func main() {
	flag.Parse()

	server := http.ServerInstance(
		http.RequestHandler(
			http.Routes(app.AppConfig{BaseDirectory: *baseDirectory}),
			http.LoggingMiddleware,
		),
	)

	server.ReadServeAddr()
	server.ListenAndServe()
	server.Watch()

	server.LogErrors()
}
