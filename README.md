# A composite golang microservice that uses eventbus (pub/sub), gRPC with protobuf and rest API's

## This is a microservice that 
- uses an eventbus for pub/sub from server to client (without the need of a broker service i.e kafka, redis amq) etc
  - forked from https://github.com/asaskevich/EventBus
- uses gRPC to handle all communication between service (client and server)
- use protobuf for payloads
- has a rest API for web/cli interaction

## Envars used (example on remote server)
```bash
CLIENT_ADDRESS=192.168.1.14:7020
CLIENT_PATH=_test_
GRPCSERVER_ADDRESS=192.168.1.3:7030
LOG_LEVEL=trace
NAME=golang-eventbus-grpc
RPCSERVER_ADDRESS=192.168.1.3
RPCSERVER_PORT=:7000
SERVER_PORT=9000
VERSION=1.0.1
```

## Pipeline runs
- Update for tekton pipline run 2

