package tcp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model"
	"github.com/hellgate75/go-network/tcp/stream"
	"net"
	"sync"
	"time"
)

type Signal byte

const (
	shutdown	Signal = iota + 1
	purge
	exit
)

var (
	ServerWaitTimeout = 120 * time.Second
)
type tcpServer struct {
	sync.Mutex
	config			*model.TcpServerConfig
	running			bool
	internal		chan Signal
	commands		chan Signal
	logger			log.Logger
	activeRequests	int64
	handlers		[]*model.TcpCallHandler
	tcpListener		*net.Listener
	timer			*time.Ticker
	serverMap		map[string]interface{}
}

func(server *tcpServer) Init(config model.TcpServerConfig) (model.TcpServer, error) {
	if server.running {
		return server, errors.New(fmt.Sprint("TcpServer.Init() - Server is still running"))
	}
	server.config = &config
	return server, nil
}

func(server *tcpServer) Start() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("TcpServer.Start() - Error: %v", r))
			server.logger.Fatalf("TcpServer.Start() -  Error: %v", err)
		}
	}()
	if server.running {
		server.logger.Fatal("TcpServer.Start() - Error: Server already running")
		return errors.New(fmt.Sprint("TcpServer.Start() - Error: Server already running"))
	}
	if server.config == nil {
		server.logger.Fatal("TcpServer.Start() - Error: No server configuration provided")
		return errors.New(fmt.Sprint("TcpServer.Start() - Error: No server configuration provided"))
	}
	server.internal = make(chan Signal)
	server.commands = make(chan Signal)
	var address = fmt.Sprintf("%s:%v", server.config.Host, server.config.Port)
	if server.config.Port <= 0 {
		address = fmt.Sprintf("%s", server.config.Host)
	}
	var l net.Listener
	if server.config.Config != nil {
		l, err = tls.Listen(server.config.Network, address, server.config.Config)
	} else {
		l, err = net.Listen(server.config.Network, address)
	}
	if err == nil {
		server.running = true
		server.logger.Infof("TcpServer.Start() - Server started on: %s", address)
		server.tcpListener = &l
		go server.acceptClients()
	} else {
		server.logger.Errorf("TcpServer.Start() - Server failed to start on: %s, due to error: %v", address, err)
		server.tcpListener = nil
	}
	return err
}

func (server *tcpServer) handleConnection(conn net.Conn) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("TcpServer.handleConnection() - Error: %v", r))
			server.logger.Fatalf("TcpServer.handleConnection() - Error: %v", err)
		}
	}()
	addr := conn.RemoteAddr()
	if len(server.handlers) > 0 {
		defer func() {
			server.logger.Debugf("TcpServer.handleConnection() - Closing connection with address %+v...", addr)
			err = conn.Close()
			if err != nil {
				server.logger.Fatalf("TcpServer.handleConnection() - Close connection with address %+v - Error: %v", addr, err)
			}
		}()
		rwCloser := stream.NewConnReaderWriterCloser()
		rwCloser.Enroll(conn)
		var wg = sync.WaitGroup{}
		for _, handler := range server.handlers{
			if handler != nil {
				go func(connection net.Conn, rw stream.ConnReaderWriterCloser) {
					wg.Add(1)
					server.register()
					server.logger.Debugf("Handling request from %+v to handler named: %s", addr, (*handler).GetName())
					(*handler).HandleRequest(connection, rw)
					server.deregister()
					wg.Done()
				}(conn, rwCloser)
			}
		}
		time.Sleep(1 * time.Second)
		wg.Wait()
		_ = rwCloser.Close()
	} else {
		server.logger.Fatalf("TcpServer.acceptClients() - Closing connection from %+v for no handlers ...", conn.RemoteAddr())
		err = conn.Close()
		if err != nil {
			server.logger.Fatalf("TcpServer.handleConnection() - Close connection for address %+v - Error: %v", addr, err)
		}
	}
}

func (server *tcpServer) acceptClients() {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("TcpServer.acceptClients() - Error: %v", r))
			server.logger.Fatalf("TcpServer.acceptClients() - Error: %v", err)
		}
	}()
	if server.tcpListener != nil {
		for server.running{
			var conn net.Conn
			conn, err = (*server.tcpListener).Accept()
			if err != nil{
				server.logger.Errorf("TcpServer.acceptClients() - Acceptance Error: %v", err)
				continue
			}
			server.logger.Debugf("TcpServer.acceptClients() - Handling request from: %+v ...", conn.RemoteAddr())
			go server.handleConnection(conn)
		}
	} else {
		server.logger.Fatalf("TcpServer.acceptClients() - Invalid listener - Stopping server ...")
		err = server.Stop()
		server.logger.Error(err)
	}
}

func (server *tcpServer) shutdownTimer() {
	server.timer = time.NewTicker(5 * time.Second)
	defer func() {
		server.timer.Stop()
		server.timer = nil
	}()
	for {
		select {
		case <- server.timer.C:
			if server.checkExit() {
				server.evacuate()
				return
			}

		}
	}
}


func (server *tcpServer) register() {
	defer server.Unlock()
	server.Lock()
	server.activeRequests++
}

func (server *tcpServer) deregister() {
	defer server.Unlock()
	server.Lock()
	server.activeRequests--
}

func (server *tcpServer) checkExit() bool {
	if ! server.Running() && ! server.Working() {
		server.internal <- exit
		return true
	}
	return false
}


func (server *tcpServer) evacuate() {
	close(server.internal)
	server.internal = nil
	close(server.commands)
	server.commands = nil
}

func(server *tcpServer) Stop() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("TcpServer.Stop() - Error: %v", r))
			server.logger.Fatalf("TcpServer.Stop() -  Error: %v", err)
		}
	}()
	if ! server.running {
		err = errors.New(fmt.Sprint("TcpServer.Stop() - Error: Server is already stopped"))
		server.logger.Errorf("TcpServer.Stop() -  %v", err)
		return err
	}
	server.internal <- shutdown
	go server.shutdownTimer()
	server.running = false
	if server.tcpListener != nil {
		err := (*server.tcpListener).Close()
		if err != nil {
			server.logger.Errorf("TcpServer.Stop() - Gently shutting down server error occurred: %v", err)
			server.logger.Warnf("TcpServer.Stop() - Try brute-force server close ...")
			server.activeRequests = 0
		}
		server.tcpListener = nil
	}
	return err
}

func(server *tcpServer) Running() bool {
	return server.running
}

func(server *tcpServer) Working() bool {
	return server.activeRequests > 0
}

func(server *tcpServer) Wait() {
	defer func() {
		if r := recover(); r != nil {
			server.logger.Fatalf("TcpServer.Wait() - Fatal Error: %v", r)
		}
	}()
	server.logger.Debugf("TcpServer.Wait() - Waiting for server shutdown")
waitCycle:
	for server.Running() || server.Working() {
		select {
		case sig := <- server.internal:
			if sig == shutdown {
				server.commands	<- purge
			} else if sig == exit {
				break waitCycle
			}
		case <- time.After(ServerWaitTimeout):
			continue
		}
	}
	server.logger.Warnf("TcpServer.Wait() - Server shutdown in progress, exiting ...")
	time.Sleep(10 * time.Second)
	server.logger.Warnf("TcpServer.Wait() - exit")
}
func(server *tcpServer) containsHandler(name string) bool {
	for _, handler := range server.handlers {
		if (*handler).GetName() == name {
			return true
		}
	}
	return false
}

func(server *tcpServer) AddPath(handler model.TcpCallHandler) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("TcpServer.AddPath() - Error: %v", r))
			server.logger.Fatalf("TcpServer.AddPath() -  Error: %v", err)
		}
	}()
	var name= handler.GetName()
	if len(handler.Names()) == 0 {
		err = errors.New(fmt.Sprint("Provided handler has not method implementation"))
		server.logger.Warnf("TcpServer.AddPath() - No Web Methods for Tcp handler with name: %s", name)
	} else if server.containsHandler(name) {
		err = errors.New(fmt.Sprintf("Provided handler has duplicsted handler with name: %s", name))
		server.logger.Warnf("TcpServer.AddPath() - Duplicated Tcp handler with name: %s", name)
	} else {
		server.logger.Debugf("TcpServer.AddPath() - Duplicated Tcp handler with name: %s", name)
		handler.SetServerMap(&server.serverMap)
		handler.SetLogger(server.logger)
		handler.SetEncoding(server.config.Encoding)
		server.handlers = append(server.handlers, &handler)
		server.logger.Debugf("TcpServer.AddPath() - Adding Tcp handler with name: %s", name)
	}
	return err
}

func NewTcpServer(appName string, verbosity log.LogLevel) model.TcpServer {
	return &tcpServer{
		config: nil,
		running: false,
		logger: log.NewLogger(appName, verbosity),
		handlers: make([]*model.TcpCallHandler, 0),
		tcpListener: nil,
		timer: nil,
		serverMap: make(map[string]interface{}),
	}
}