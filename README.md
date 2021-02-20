# A composite golang microservice that uses eventbus (pub/sub), gRPC with protobuf and rest API's

## This is a microservice that 
- uses an eventbus for pub/sub from server to client (without the need of a broker service i.e kafka, redis amq) etc
- uses gRPC to handle all communication between service (client and server)
- use protobuf for payloads
- has a rest API for web/cli interaction
