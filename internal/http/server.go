package http

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Server struct {
	httpServer *http.Server
	err        error

	errChannel chan error
	sigChannel chan os.Signal
}

var (
	once   sync.Once
	server *Server
)

const (
	shutdownWaitDuration = time.Second * 15
	writeTimeout         = time.Second * 15
	readTimeout          = time.Second * 15
	idleTimeout          = time.Second * 60
)

func interruptChannel() chan os.Signal {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	return channel
}

func ServerInstance(handler http.Handler) *Server {
	once.Do(func() {
		server = &Server{
			httpServer: &http.Server{
				WriteTimeout: writeTimeout,
				ReadTimeout:  readTimeout,
				IdleTimeout:  idleTimeout,
				Handler:      handler,
			},

			errChannel: make(chan error, 1),
			sigChannel: interruptChannel(),
		}
	})

	return server
}

func (server *Server) LogErrors() {
	if server.err != nil {
		log.Println(server.err)
	}
}

func (server *Server) obtainPort() string {
	if server.err != nil {
		return "0"
	}

	port := 0

	var listener net.Listener
	listener, server.err = net.Listen("tcp", ":0")

	if server.err == nil {
		port = listener.
			Addr().(*net.TCPAddr).Port

		server.err = listener.Close()
	}

	return fmt.Sprint(port)
}

func (server *Server) ReadServeAddr() {
	port := server.obtainPort()

	if server.err != nil {
		return
	}

	server.httpServer.Addr =
		fmt.Sprintf(":%s", port)

	log.Printf(
		"Listening on address: %s",
		server.httpServer.Addr)
}

func (server *Server) ListenAndServe() {
	if server.err != nil {
		return
	}

	go func() {
		server.errChannel <- server.
			httpServer.ListenAndServe()
	}()
}

func (server *Server) shutdown() {
	waitContext, cancel := context.WithTimeout(
		context.Background(),
		shutdownWaitDuration,
	)
	defer cancel()

	log.Println("Shutting down...")

	server.httpServer.Shutdown(waitContext)
}

func (server *Server) Watch() {
	if server.err != nil {
		return
	}

	select {
	case server.err = <-server.errChannel:
	case <-server.sigChannel:
		log.Println("Received ^C.")
		server.shutdown()
	}
}
