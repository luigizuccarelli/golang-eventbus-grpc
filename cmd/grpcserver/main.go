// +build real

package main

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/impl"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/service"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/validator"
	"github.com/microlib/simple"
	gRPC "google.golang.org/grpc"
)

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
	port, _ := strconv.Atoi(os.Getenv("GRPCSERVER_PORT"))
	l := getNetListener(port, logger)
	gRPCServer := gRPC.NewServer()
	dataschemaServiceImpl := impl.NewDataSchemaServiceGrpcImpl()
	service.RegisterDataSchemaServiceServer(gRPCServer, dataschemaServiceImpl)

	// probes setup
	healthService := impl.NewHealthServiceGrpcImpl()
	service.RegisterHealthServer(gRPCServer, healthService)

	logger.Info(fmt.Sprintf("GRPC Server starting on port %d", port))
	if err := gRPCServer.Serve(l); err != nil {
		logger.Error(fmt.Sprintf("failed to serve: %v", err))
		os.Exit(-1)
	}

}

func getNetListener(port int, logger *simple.Logger) net.Listener {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to listen: %v", err))
		os.Exit(-1)
	}
	return lis
}
