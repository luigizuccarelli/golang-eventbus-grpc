// +build fake

package connectors

import (
	"fmt"

	"github.com/microlib/simple"
)

type FakeRPCClient struct {
}

func (f *FakeRPCClient) Close() {
}

func (f *FakeRPCClient) Call(regservice string, args interface{}, flag *bool) error {
	return nil
}

type RPCClient struct {
	Connector *FakeRPCClient
}

// Fake all connections
type Connectors struct {
	Logger *simple.Logger
}

func (c Connectors) Error(msg string, val ...interface{}) {
	c.Logger.Error(fmt.Sprintf(msg, val...))
}

func (c Connectors) Info(msg string, val ...interface{}) {
	c.Logger.Info(fmt.Sprintf(msg, val...))
}

func (c Connectors) Debug(msg string, val ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (c Connectors) Trace(msg string, val ...interface{}) {
	c.Logger.Trace(fmt.Sprintf(msg, val...))
}

func (c Connectors) DialHttpPath(addr string, path string) (*RPCClient, error) {
	return &RPCClient{Connector: &FakeRPCClient{}}, nil
}

func NewTestConnectors(logger *simple.Logger) Clients {
	conns := &Connectors{Logger: logger}
	return conns
}
