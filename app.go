package main

type AppConfig struct {
	BaseDirectory string
}

func AppServer(config AppConfig) *Server {
	return ServerInstance(
		RequestHandler(Routes(config),
			LoggingMiddleware,
		),
	)
}

func ServeApp(config AppConfig) {
	server := AppServer(config)

	server.ReadServeAddr()
	server.ListenAndServe()
	server.Watch()

	server.LogErrors()
}
