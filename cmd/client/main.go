// +build real

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	// go clean -modcache (for version problems)
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/connectors"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/domain"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/eventbus"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/service"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/validator"
	"github.com/microlib/simple"
	gRPC "google.golang.org/grpc"
)

// webhook private function - this is the callback function
func webhook(conn *connectors.Connectors) {
	conn.Info("Function webhook received message from producer")
	s := os.Getenv("GRPCSERVER_ADDRESS")
	grpcConn, e := gRPC.Dial(s, gRPC.WithInsecure())
	if e != nil {
		conn.Error("Function webhook response from server %v", e)
		return
	}
	defer grpcConn.Close()
	c := service.NewDataSchemaServiceClient(grpcConn)
	dataschemaModel := domain.DataSchema{
		Id:      int64(456),
		Name:    string("Grpc-Demo"),
		Status:  string("OK"),
		Payload: string("{\"message\":\"dude this sh*t is working\"}"),
	}
	if responseMessage, e := c.Get(context.Background(), &dataschemaModel); e != nil {
		conn.Error("Function webhook response from server %v", e)
		return
	} else {
		conn.Info("DataSchema from server (GET rpc)")
		conn.Info("Data : %v", responseMessage)
	}
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

	bus := eventbus.New()
	client := eventbus.NewClient(os.Getenv("CLIENT_ADDRESS"), os.Getenv("CLIENT_PATH"), bus)
	e := client.Start(conn)
	if e != nil {
		conn.Error("Client : %v", e)
	}
	client.Subscribe(os.Getenv("CLIENT_PATH"), webhook, os.Getenv("RPCSERVER_ADDRESS")+os.Getenv("RPCSERVER_PORT"), "/_server_bus_", conn)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	exit_chan := make(chan int)

	go func() {
		for {
			s := <-c
			switch s {
			case syscall.SIGHUP:
				exit_chan <- 0
			case syscall.SIGINT:
				exit_chan <- 0
			case syscall.SIGTERM:
				exit_chan <- 0
			case syscall.SIGQUIT:
				exit_chan <- 0
			default:
				exit_chan <- 1
			}
		}
	}()

	code := <-exit_chan
	client.Unsubscribe(os.Getenv("CLIENT_TOPIC"), webhook)
	client.Stop()
	fmt.Println("client shutdown successfully")
	os.Exit(code)
}
