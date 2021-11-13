package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"

	"github.com/arindas/elevate/pkg/app"
	"github.com/arindas/elevate/pkg/http"
)

var baseDirectory = flag.String("base_dir", ".", "Base directory for storing files.")

//go:embed web/*
var embeddedContent embed.FS

const webDirectory = "web"

func AppConfig() app.AppConfig {
	staticContent, err := fs.Sub(embeddedContent, webDirectory)
	if err != nil {
		log.Printf("Unable to embed static files.")
	}

	flag.Parse()
	log.Printf("Using base directory: %s", *baseDirectory)

	return app.AppConfig{
		BaseDirectory: *baseDirectory,
		StaticContent: staticContent,
	}
}

func main() {
	server := http.ServerInstance(
		http.RequestHandler(
			http.Routes(AppConfig()),
			http.LoggingMiddleware,
		),
	)

	server.ReadServeAddr()
	server.ListenAndServe()
	server.Watch()

	server.LogErrors()
}
