package main

import (
	"flag"
	"log"

	"github.com/arindas/elevate/internal/app"
	"github.com/arindas/elevate/internal/http"
)

var baseDirectory = flag.String("base_dir", ".", "Base directory for storing files.")

func main() {
	flag.Parse()

	log.Printf("Base directory: %s", *baseDirectory)

	server := http.ServerInstance(
		http.RequestHandler(
			http.Routes(app.AppConfig{
				BaseDirectory: *baseDirectory,
			}),
			http.LoggingMiddleware,
		),
	)

	server.ReadServeAddr()
	server.ListenAndServe()
	server.Watch()

	server.LogErrors()
}
