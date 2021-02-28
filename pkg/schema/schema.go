package schema

import "github.com/luigizuccarelli/golang-eventbus-grpc/pkg/connectors"

type Request struct {
	Id      string `json:"id,omitempty"`
	Message string `json:"message"`
}

// Response schema
type Response struct {
	Code    int    `json:"code,omitempty"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// SubscribeType - how the client intends to subscribe
type SubscribeType int

// SubscribeArg - object to hold subscribe arguments from remote event handlers
type SubscribeArg struct {
	ClientAddr    string
	ClientPath    string
	ServiceMethod string
	SubscribeType SubscribeType
	Topic         string
}

// ClientArg - object containing event for client to publish locally
type ClientArg struct {
	Args  *connectors.Connectors
	Topic string
}
