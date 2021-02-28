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
	Subscribe schema.SubscribeType = iota
	SubscribeOnce
	REGISTERSERVICE = "ServerService.Register"
)

// we have to kee this local

// Server - object capable of being subscribed to by remote handlers
type Server struct {
	eventBus    Bus
	address     string
	path        string
	Subscribers map[string][]*schema.SubscribeArg
	service     *ServerService
}

// ServerService - service object to listen to remote subscriptions
type ServerService struct {
	server  *Server
	wg      *sync.WaitGroup
	started bool
}

// NewServer - create a new Server at the address and path
func NewServer(address, path string, eventBus Bus) *Server {
	server := new(Server)
	server.eventBus = eventBus
	server.address = address
	server.path = path
	server.Subscribers = make(map[string][]*schema.SubscribeArg)
	server.service = &ServerService{server, &sync.WaitGroup{}, false}
	return server
}

// EventBus - returns wrapped event bus
func (server *Server) EventBus() Bus {
	return server.eventBus
}

// rpcCallBack - the callback function must pass the connectors.Connector parameter
func (server *Server) rpcCallback(subscribeArg *schema.SubscribeArg) func(conn *connectors.Connectors) {
	return func(conn *connectors.Connectors) {
		client, connErr := conn.DialHttpPath(subscribeArg.ClientAddr, subscribeArg.ClientPath)
		if connErr != nil {
			conn.Error("rpcCallback client connection: %v", connErr)
			delete(server.Subscribers, subscribeArg.Topic)
			conn.Info("rpcCallback subscriber removed: %s", subscribeArg.Topic)
			// call the Handler to cleanup
			hnd, _ := server.eventBus.Handler(subscribeArg.ClientPath, subscribeArg.Topic)
			conn.Info("rpcCallback subscriber remaining: %v", server.Subscribers)
			conn.Info("rpcCallback handler remaining: %v", hnd)
			return
		}
		defer client.Connector.Close()
		clientArg := new(schema.ClientArg)
		clientArg.Topic = subscribeArg.Topic
		clientArg.Args = conn
		var reply bool
		err := client.Connector.Call(subscribeArg.ServiceMethod, clientArg, &reply)
		if err != nil {
			conn.Error("rcpCallback call : %v", err)
		}
	}
}

// HasClientSubscribed - True if a client subscribed to this server with the same topic
func (server *Server) HasClientSubscribed(arg *schema.SubscribeArg) bool {
	if topicSubscribers, ok := server.Subscribers[arg.Topic]; ok {
		for _, topicSubscriber := range topicSubscribers {
			if *topicSubscriber == *arg {
				return true
			}
		}
	}
	return false
}

// Start - starts a service for remote clients to subscribe to events
func (server *Server) Start(conn connectors.Clients) error {
	var err error
	service := server.service
	if !service.started {
		rpcServer := rpc.NewServer()
		rpcServer.Register(service)
		rpcServer.HandleHTTP(server.path, "/debug"+server.path)
		l, e := net.Listen("tcp", server.address)
		if e != nil {
			err = e
			conn.Error("Function Start : listen error: %v", e)
		}
		service.started = true
		service.wg.Add(1)
		go http.Serve(l, nil)
		conn.Info("Function Start : rpc server started %s", server.address)
	} else {
		err = errors.New("server bus already started")
	}
	return err
}

// Stop - signal for the service to stop serving
func (server *Server) Stop() {
	service := server.service
	if service.started {
		service.wg.Done()
		service.started = false
	}
}

// Register - Registers a remote handler to this event bus
// for a remote subscribe - a given client address only needs to subscribe once
// event will be republished in local event bus
func (service *ServerService) Register(arg *schema.SubscribeArg, success *bool) error {
	subscribers := service.server.Subscribers
	if !service.server.HasClientSubscribed(arg) {
		rpcCallback := service.server.rpcCallback(arg)
		switch arg.SubscribeType {
		case Subscribe:
			service.server.eventBus.Subscribe(arg.ClientPath, arg.Topic, rpcCallback)
		case SubscribeOnce:
			service.server.eventBus.SubscribeOnce(arg.Topic, rpcCallback)
		}
		var topicSubscribers []*schema.SubscribeArg
		if _, ok := subscribers[arg.Topic]; ok {
			topicSubscribers = []*schema.SubscribeArg{arg}
		} else {
			topicSubscribers = subscribers[arg.Topic]
			topicSubscribers = append(topicSubscribers, arg)
		}
		subscribers[arg.Topic] = topicSubscribers
	}
	*success = true
	return nil
}
