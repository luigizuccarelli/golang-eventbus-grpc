# A composite golang microservice that uses eventbus (pub/sub), gRPC with protobuf and rest API's

![Quality Gate Status](status-badge.png)  ![Code Coverage](code-coverage.png)

## This is a microservice that 
- uses an eventbus for pub/sub from server to client (without the need of a broker service i.e kafka, redis amq) etc
  - forked from https://github.com/asaskevich/EventBus
- uses gRPC to handle all communication between service (client and server)
- use protobuf for payloads
- has a rest API for web/cli interaction

## Envars used (example on remote server)
```bash
export CLIENT_ADDRESS=192.168.1.3:7020
export CLIENT_PATH=_test_
export GRPCSERVER_ADDRESS=192.168.1.3:7030
export LOG_LEVEL=trace
export NAME=golang-eventbus-grpc
export RPCSERVER_ADDRESS=192.168.1.3
export RPCSERVER_PORT=:7000
export SERVER_PORT=9000
export VERSION=1.0.1
```

## Pipeline runs
- Update for tekton pipline run 11


## Curl API service

```
curl -d'{"id":"20","message":"hello reedge"}' http://127.0.0.1:9000/api/v1/service
```
