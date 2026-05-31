package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"vds.io/goml/siren"
)

var ErrWorkInProgress = errors.New("work in progress")

type Server struct {
	httpServer     http.Server
	workInProgress bool
	listenResult   chan error

	alarm    *siren.Siren
	schedule []TimeOfDay
}

func NewServer(alarm *siren.Siren) *Server {
	s := Server{
		httpServer: http.Server{
			Addr: ":13086",
		},
		listenResult: make(chan error),
		alarm:        alarm,
	}
	s.httpServer.Handler = &s

	return &s
}

func (server *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	r, _ := server.alarm.Start()
	go func() {
		err := <-r
		fmt.Println(err)
	}()

	fmt.Fprintf(resp, "Get Off My Lawn!")
}

func (server *Server) startListening() {
	server.listenResult <- server.httpServer.ListenAndServe()
	close(server.listenResult)
}

func (server *Server) StartListening() (chan error, error) {
	if server.workInProgress {
		return nil, ErrWorkInProgress
	}

	go server.startListening()

	return server.listenResult, nil
}

func (server *Server) Stop(ctx context.Context) error {
	return server.httpServer.Shutdown(ctx)
}
