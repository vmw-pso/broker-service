package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/vmw-pso/toolkit"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)

	var (
		port = flags.Int("port", 80, "port to listen on")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	addr := fmt.Sprintf(":%d", *port)
	srv := newServer()

	return http.ListenAndServe(addr, srv)
}

type server struct {
	mux   *chi.Mux
	tools toolkit.Tools
}

func newServer() *server {
	mux := chi.NewMux()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	tools := toolkit.Tools{}

	srv := &server{
		mux:   mux,
		tools: tools,
	}
	srv.routes()

	return srv
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *server) routes() {
	s.mux.Post("/handle", s.handleRequest())
}
