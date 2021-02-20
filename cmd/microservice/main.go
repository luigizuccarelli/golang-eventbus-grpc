// +build real

package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/connectors"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/eventbus"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/handlers"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/validator"
	"github.com/microlib/simple"
)

// startServers private function that starts all relevant servers
func startServers(conn connectors.Clients) *http.Server {
	// start the rpc server first
	rpcServer := eventbus.NewServer(os.Getenv("RPCSERVER_PORT"), "/_server_bus_", eventbus.New())
	rpcServer.Start(conn)

	// start the rest API server
	srv := &http.Server{Addr: ":" + os.Getenv("SERVER_PORT")}
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/service", func(w http.ResponseWriter, req *http.Request) {
		handlers.ServiceHandler(w, req, conn, rpcServer)
	}).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/v1/sys/info/isalive", handlers.IsAlive).Methods("GET")
	http.Handle("/", r)
	if err := srv.ListenAndServe(); err != nil {
		conn.Error("Httpserver: ListenAndServe() error: " + err.Error())
	}
	return srv
}

func main() {
	var logger *simple.Logger
	if os.Getenv("LOG_LEVEL") == "" {
		logger = &simple.Logger{Level: "info"}
	} else {
		logger = &simple.Logger{Level: os.Getenv("LOG_LEVEL")}
	}
	err := validator.ValidateEnvars(logger)
	if err != nil {
		os.Exit(-1)
	}
	conn := connectors.NewClientConnections(logger)
	srv := startServers(conn)
	conn.Info("Starting server on port " + srv.Addr)
}
