// +build fake

package eventbus

import (
	"testing"

	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/connectors"
	"github.com/microlib/simple"
)

func TestNewClient(t *testing.T) {
	logger := &simple.Logger{Level: "trace"}
	bus := New()
	if bus == nil {
		t.Log("New EventBus not created!")
		t.Fail()
	}

	client := NewClient("test", "path", bus)
	if client == nil {
		t.Log("New client not created!")
		t.Fail()
	}

	// lets get the underlying eventbus
	b := client.EventBus()
	if b == nil {
		t.Log("Client eventbus method failed!")
		t.Fail()
	}

	conn := connectors.NewTestConnectors(logger)
	handler := func() {}
	// lets subscribe
	client.Subscribe("topic", handler, "serveraddr", "servepath", conn)
	client.SubscribeOnce("topic", handler, "serveraddr", "servepath", conn)
	client.Unsubscribe("topic", handler)
}

func TestClientService(t *testing.T) {
	logger := &simple.Logger{Level: "trace"}
	bus := New()
	client := NewClient("test", "path", bus)
	conn := connectors.NewTestConnectors(logger)
	err := client.Start(conn)
	if err != nil {
		t.Fail()
	}
	client.Stop()

	client = NewClient("newtest", "newpath", bus)
	client.address = ":10000"
	err = client.Start(conn)
	if err != nil {
		t.Fail()
	}
	client.Stop()
}
