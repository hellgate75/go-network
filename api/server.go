package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model"
	"net/http"
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

type apiServer struct {
	sync.Mutex
	config			*model.ServerConfig
	running			bool
	router			*mux.Router
	internal		chan Signal
	commands		chan Signal
	logger			log.Logger
	activeRequests	int64
	handlers		map[string]*model.ApiCallHandler
	httpServer		*http.Server
	timer			*time.Ticker
	serverMap		map[string]interface{}
}

func (server *apiServer) Init(config model.ServerConfig) (model.ApiServer, error) {
	if server.running {
		return server, errors.New(fmt.Sprint("Server is still running"))
	}
	server.config = &config
	return server, nil
}

func (server *apiServer) Start() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("ApiServer.Start() - Error: %v", r))
			server.logger.Fatalf("%v", err)
		}
	}()
	if server.running {
		server.logger.Fatal("ApiServer.Start() - Error: Server already running")
		return errors.New(fmt.Sprint("ApiServer.Start() - Error: Server already running"))
	}
	if server.config == nil {
		server.logger.Fatal("ApiServer.Start() - Error: No server configuration provided")
		return errors.New(fmt.Sprint("ApiServer.Start() - Error: No server configuration provided"))
	}
	var address = fmt.Sprintf("%s:%v", server.config.Host, server.config.Port)
	server.httpServer = &http.Server{
		Addr: address,
		Handler: server.router,
		TLSConfig: server.config.Config,
	}
	if server.config.CertPath != "" && server.config.KeyPath != "" {
		// TLS encryption
		server.logger.Debugf("ApiServer.Start() - Running TLS encryption listener on: %s", address)
		err = server.httpServer.ListenAndServeTLS(server.config.CertPath, server.config.KeyPath)

	} else {
		// No TLS encryption
		server.logger.Debugf("ApiServer.Start() - Running non-TLS encryption listener on: %s", address)
		err = server.httpServer.ListenAndServe()
	}
	if err == nil {
		server.running = true
		server.logger.Infof("ApiServer.Start() - Server started on: %s", address)
	} else {
		server.logger.Errorf("ApiServer.Start() - Server failed to start on: %s, due to error: %v", address, err)
	}
	return err
}

func (server *apiServer) Stop() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("ApiServer.Stop() - Error: %v", r))
			server.logger.Fatalf("%v", err)
		}
	}()
	if ! server.running {
		err = errors.New(fmt.Sprint("ApiServer.Stop() - Error: Server is already stopped"))
		server.logger.Errorf("%v", err)
		return err
	}
	server.internal <- shutdown
	go server.shutdownTimer()
	server.running = false
	if server.httpServer != nil {
		err := server.httpServer.Shutdown(context.Background())
		if err != nil {
			server.logger.Errorf("ApiServer.evacuate() - Gently shutting down server error occurred: %v", err)
			server.logger.Warnf("ApiServer.evacuate() - Try brute-force server close ...")
			err = server.httpServer.Close()
			if err != nil {
				server.logger.Errorf("ApiServer.evacuate() - Brute close server error occurred: %v", err)
			}
			server.activeRequests = 0
		}
		server.httpServer = nil
	}
	return err
}

func (server *apiServer) Running() bool {
	return server.running
}

func (server *apiServer) register() {
	defer server.Unlock()
	server.Lock()
	server.activeRequests++
}

func (server *apiServer) deregister() {
	defer server.Unlock()
	server.Lock()
	server.activeRequests--
}

func (server *apiServer) checkExit() bool {
	if ! server.Running() && ! server.Working() {
		server.internal <- exit
		return true
	}
	return false
}

func (server *apiServer) shutdownTimer() {
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

func (server *apiServer) evacuate() {
	close(server.internal)
	close(server.commands)
}

func (server *apiServer) Working() bool {
	return server.activeRequests > 0
}

func (server *apiServer) Wait() {
	defer func() {
		if r := recover(); r != nil {
			server.logger.Fatalf("ApiServer.Wait() - Fatal Error: %v", r)
		}
	}()
	server.logger.Debugf("ApiServer.Wait() - Waiting for server shutdown")
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
	server.logger.Warnf("ApiServer.Wait() - Server shutdown in progress, exiting ...")
	time.Sleep(10 * time.Second)
	server.logger.Warnf("ApiServer.Wait() - exit")
}

func (server *apiServer) AddPath(handler model.ApiCallHandler) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("ApiServer.AddPath() - Error: %v", r))
			server.logger.Fatalf("%v", err)
		}
	}()
	var path = handler.GetPath()
	if len(path) == 0 {
		err = errors.New(fmt.Sprint("Provided handler has empty path"))
		server.logger.Warn("ApiServer.AddPath() - Empty Path for Api handler")
	} else if len(handler.Methods()) == 0 {
		err = errors.New(fmt.Sprint("Provided handler has not method implementation"))
		server.logger.Warnf("ApiServer.AddPath() - No Web Methods for Api handler in path %s", path)
	} else if _, ok := server.handlers[path]; ok {
		err = errors.New(fmt.Sprintf("Provided handler has duplicsted path: %s", path))
		server.logger.Warnf("ApiServer.AddPath() - Duplicated Api handler for path %s", path)
	} else {
		handler.SetServerMap(&server.serverMap)
		handler.SetLogger(server.logger)
		server.router.HandleFunc(path, handler.HandleRequest).Methods(handler.Methods()...)
		server.handlers[path]=&handler
		server.logger.Debugf("ApiServer.AddPath() - Adding Api handler for path %s", path)
	}
	return err
}

func NewApiServer(appName string, verbosity log.LogLevel) model.ApiServer {
	return &apiServer{
		config: nil,
		running: false,
		router: mux.NewRouter(),
		internal: make(chan Signal),
		commands: make(chan Signal),
		logger: log.NewLogger(appName, verbosity),
		handlers: make(map[string]*model.ApiCallHandler),
		httpServer: nil,
		timer: nil,
		serverMap: make(map[string]interface{}),
	}
}