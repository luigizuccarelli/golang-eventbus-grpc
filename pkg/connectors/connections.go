// +build real

package connectors

import (
	"fmt"
	"net/rpc"

	"github.com/microlib/simple"
)

// Connections struct - all backend connections in a common object
type Connectors struct {
	Logger *simple.Logger
}

type RPCClient struct {
	Connector *rpc.Client
}

// NewClientConnections - function that initialises all client (3rd party) interfaces
func NewClientConnections(logger *simple.Logger) Clients {
	return &Connectors{Logger: logger}
}

func (c *Connectors) Error(msg string, val ...interface{}) {
	c.Logger.Error(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Info(msg string, val ...interface{}) {
	c.Logger.Info(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (c *Connectors) Trace(msg string, val ...interface{}) {
	c.Logger.Trace(fmt.Sprintf(msg, val...))
}

// DialHttPath - wrapper for the actual rcp DialHTTPPath call. This allows us to use a fake/mock for testing
func (c *Connectors) DialHttpPath(addr string, path string) (*RPCClient, error) {
	client, err := rpc.DialHTTPPath("tcp", addr, path)
	if err != nil {
		return &RPCClient{}, err
	}
	return &RPCClient{Connector: client}, nil
}
