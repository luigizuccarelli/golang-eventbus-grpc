package eventbus

import (
	"errors"
	"net"
	"net/http"
	"net/rpc"
	"sync"

	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/connectors"
	"github.com/luigizuccarelli/golang-eventbus-grpc/pkg/schema"
)

const (
	// PublishService - Client service method
	PUBLISHSERVICE = "ClientService.PushEvent"
)

// we have t make this local

// Client - object capable of subscribing to a remote event bus
type Client struct {
	eventBus Bus
	address  string
	path     string
	service  *ClientService
}

// ClientService - service object listening to events published in a remote event bus
type ClientService struct {
	client  *Client
	wg      *sync.WaitGroup
	started bool
}

// NewClient - create a client object with the address and server path
func NewClient(address, path string, eventBus Bus) *Client {
	client := new(Client)
	client.eventBus = eventBus
	client.address = address
	client.path = path
	client.service = &ClientService{client, &sync.WaitGroup{}, false}
	return client
}

// EventBus - returns the underlying event bus
func (client *Client) EventBus() Bus {
	return client.eventBus
}

// added lmz
func (client *Client) Unsubscribe(topic string, fn interface{}) {
	client.eventBus.Unsubscribe(topic, fn)
}

func (client *Client) doSubscribe(topic string, fn interface{}, serverAddr, serverPath string, subscribeType schema.SubscribeType, conn connectors.Clients) {
	defer func() {
		if r := recover(); r != nil {
			conn.Error("Server not found: %v", r)
		}
	}()

	rpcClient, err := conn.DialHttpPath(serverAddr, serverPath)
	defer rpcClient.Connector.Close()
	if err != nil {
		conn.Error("dialing: %v", err)
	}
	args := &schema.SubscribeArg{client.address, client.path, PUBLISHSERVICE, subscribeType, topic}
	reply := new(bool)
	err = rpcClient.Connector.Call(REGISTERSERVICE, args, reply)
	if err != nil {
		conn.Error("Register error: %v", err)
	}
	if *reply {
		client.eventBus.Subscribe(client.path, topic, fn)
	}
}

//Subscribe subscribes to a topic in a remote event bus
func (client *Client) Subscribe(topic string, fn interface{}, serverAddr, serverPath string, conn connectors.Clients) {
	client.doSubscribe(topic, fn, serverAddr, serverPath, Subscribe, conn)
}

//SubscribeOnce subscribes once to a topic in a remote event bus
func (client *Client) SubscribeOnce(topic string, fn interface{}, serverAddr, serverPath string, conn connectors.Clients) {
	client.doSubscribe(topic, fn, serverAddr, serverPath, SubscribeOnce, conn)
}

// Start - starts the client service to listen to remote events
func (client *Client) Start(conn connectors.Clients) error {
	var err error
	service := client.service
	if !service.started {
		server := rpc.NewServer()
		server.Register(service)
		server.HandleHTTP(client.path, "/debug"+client.path)
		conn.Info("Function Start : info %v", client)
		l, err := net.Listen("tcp", client.address)
		if err == nil {
			service.wg.Add(1)
			service.started = true
			go http.Serve(l, nil)
			conn.Info("Function Start : client server started %s", client.address)
		}
	} else {
		err = errors.New("Client service already started")
	}
	return err
}

// Stop - signal for the service to stop serving
func (client *Client) Stop() {
	service := client.service
	if service.started {
		service.wg.Done()
		service.started = false
	}
}

// PushEvent - exported service to listening to remote events
func (service *ClientService) PushEvent(arg *schema.ClientArg, reply *bool) error {
	service.client.eventBus.Publish(arg.Topic, arg.Args)
	*reply = true
	return nil
}
