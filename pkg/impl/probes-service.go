package impl

import (
	"context"
	"log"

	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/domain"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/service"
)

// HealthServiceGrpcImpl is a implementation of HealtService Grpc Service.
type HealthServiceGrpcImpl struct {
}

//NewHealthsServiceGrpcImpl returns the pointer to the implementation.
func NewHealthServiceGrpcImpl() *HealthServiceGrpcImpl {
	return &HealthServiceGrpcImpl{}
}

//Add functions implementation of gRPC Service.
func (serviceImpl *HealthServiceGrpcImpl) Check(ctx context.Context, in *domain.HealthCheckRequest) (*domain.HealthCheckResponse, error) {
	log.Println("Received health check request ")

	return &domain.HealthCheckResponse{
		Status: domain.HealthCheckResponse_SERVING,
	}, nil
}

func (serviceImpl *HealthServiceGrpcImpl) Watch(in *domain.HealthCheckRequest, server service.Health_WatchServer) error {
	log.Println("Received health watch request ")

	return server.Send(&domain.HealthCheckResponse{
		Status: domain.HealthCheckResponse_SERVING,
	})
}
