// +build fake

package eventbus

import (
	"testing"

	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/connectors"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/schema"
	"github.com/microlib/simple"
)

func TestNewServer(t *testing.T) {
	bus := New()
	if bus == nil {
		t.Log("New EventBus not created!")
		t.Fail()
	}
	server := NewServer("123", "456", bus)
	if server == nil {
		t.Log("New iServer not created!")
		t.Fail()
	}

	bus = server.EventBus()
	if bus == nil {
		t.Log("EventBus not found or null!")
		t.Fail()
	}
}

func TestHasSubscribed(t *testing.T) {
	logger := &simple.Logger{Level: "trace"}
	bus := New()
	if bus == nil {
		t.Log("New EventBus not created!")
		t.Fail()
	}
	server := NewServer("123", "456", bus)
	if server == nil {
		t.Log("New iServer not created!")
		t.Fail()
	}
	client := NewClient("test", "path", bus)
	if client == nil {
		t.Log("New client not created!")
		t.Fail()
	}

	conn := connectors.NewTestConnectors(logger)
	handler := func() {}
	// lets subscribe
	client.Subscribe("topicA", handler, "serveraddr", "servepath", conn)
	client.Subscribe("topicB", handler, "serveraddr", "servepath", conn)
	client.SubscribeOnce("topic", handler, "serveraddr", "servepath", conn)

	args := &schema.SubscribeArg{client.address, client.path, PUBLISHSERVICE, 0, "topicA"}
	b := server.HasClientSubscribed(args)
	if b {
		t.Fail()
	}
}

func TestStart(t *testing.T) {
	logger := &simple.Logger{Level: "trace"}
	bus := New()
	if bus == nil {
		t.Log("New EventBus not created!")
		t.Fail()
	}
	server := NewServer(":7000", "_tst_", bus)
	if server == nil {
		t.Log("New Server not created!")
		t.Fail()
	}

	conn := connectors.NewTestConnectors(logger)
	err := server.Start(conn)
	if err != nil {
		t.Fail()
	}
	server.Stop()
}

func TestCallBack(t *testing.T) {
	logger := &simple.Logger{Level: "trace"}
	bus := New()
	server := NewServer(":7001", "_newtest", bus)
	client := NewClient(":7002", "path", bus)
	flag := 0
	fn := func() { flag += 1 }
	conn := connectors.NewTestConnectors(logger)
	client.Subscribe("topicA", fn, ":7001", "_newtest", conn)
	args := &schema.SubscribeArg{client.address, client.path, PUBLISHSERVICE, 0, "topicA"}
	newfn := server.rpcCallback(args)
	go newfn(conn.(*connectors.Connectors))
	if args == nil {
		t.Fail()
	}
}

func TestRegister(t *testing.T) {
	bus := New()
	server := NewServer(":7001", "_newtest", bus)
	client := NewClient(":7002", "path", bus)
	args := &schema.SubscribeArg{client.address, client.path, PUBLISHSERVICE, 0, "topicA"}
	service := server.service
	var flag bool = false
	service.Register(args, &flag)
	if !flag {
		t.Fail()
	}
	flag = false
	args = &schema.SubscribeArg{client.address, client.path, PUBLISHSERVICE, 1, "topicA"}
	service.Register(args, &flag)
	if !flag {
		t.Fail()
	}
}
