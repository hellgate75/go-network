package pipe

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model"
	"io/ioutil"
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
	ServerClientResetTimeout = 500 * time.Millisecond
)

type pipeNode struct {
	config				*model.PipeNodeConfig
	running				bool
	inChan				chan model.PipeMessage
	outChan				chan model.PipeMessage
	internal			chan Signal
	commands			chan Signal
	logger				log.Logger
	timer				*time.Ticker
	activeRequests		int64
	activeClients		int64
	requestsMutex		sync.Mutex
	clientsMutex		sync.Mutex
	tcpListener			*net.Listener
	outputAddress		string
	inChanCreated		bool
	outChanCreated		bool
}

func (pipe *pipeNode) Init(config model.PipeNodeConfig) (model.PipeNode, error) {
	if pipe.running {
		return pipe, errors.New(fmt.Sprint("PipeNode.Init() - Server is still running"))
	}
	pipe.config = &config
	return pipe, nil
}

func (pipe *pipeNode) Type() model.PipeType {
	return pipe.config.Type
}

func (pipe *pipeNode) Start() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("PipeNode.Start() - Error: %v", r))
			pipe.logger.Fatalf("PipeNode.Start() -  Error: %v", err)
		}
	}()
	if pipe.running {
		pipe.logger.Fatal("PipeNode.Start() - Error: Server already running")
		return errors.New(fmt.Sprint("PipeNode.Start() - Error: Server already running"))
	}
	if pipe.config == nil {
		pipe.logger.Fatal("PipeNode.Start() - Error: No server configuration provided")
		return errors.New(fmt.Sprint("PipeNode.Start() - Error: No server configuration provided"))
	}
	if pipe.config.Type != model.InputPipe && pipe.config.Type != model.OutputPipe  && pipe.config.Type != model.InputOutputPipe {
		pipe.logger.Fatalf("PipeNode.Start() - Error: Invalid Pipe Node Type %v", pipe.config.Type)
		return errors.New(fmt.Sprintf("PipeNode.Start() - Error: Invalid Pipe Node Type %v", pipe.config.Type))
	}
	pipe.internal = make(chan Signal)
	pipe.commands = make(chan Signal)
	if pipe.config.Type == model.InputPipe || pipe.config.Type == model.InputOutputPipe {
		go func() {
			var address = fmt.Sprintf("%s:%v", pipe.config.InHost, pipe.config.InPort)
			var l net.Listener
			if pipe.config.Config != nil {
				l, err = tls.Listen("tcp", address, pipe.config.Config)
			} else {
				l, err = net.Listen("tcp", address)
			}
			if err == nil {
				pipe.running = true
				pipe.logger.Infof("PipeNode.Start() - Server started on: %s", address)
				pipe.tcpListener = &l
				go pipe.acceptClients()
			} else {
				pipe.logger.Errorf("PipeNode.Start() - Server failed to start on: %s, due to error: %v", address, err)
				pipe.tcpListener = nil
			}
		}()
	}
	if pipe.config.Type == model.OutputPipe || pipe.config.Type == model.InputOutputPipe {
		pipe.outputAddress = fmt.Sprintf("%s:%v", pipe.config.OutHost, pipe.config.OutPort)
		go pipe.readFromInputChannel()
	}
	return err
}

func (pipe *pipeNode) handleConnection(conn net.Conn) {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("PipeNode.handleConnection() - Error: %v", r))
			pipe.logger.Fatalf("PipeNode.handleConnection() - Error: %v", err)
		}
	}()
	addr := conn.RemoteAddr()
	defer func() {
		pipe.logger.Debugf("PipeNode.handleConnection() - Closing connection with address %+v...", addr)
		err = conn.Close()
		if err != nil {
			pipe.logger.Fatalf("PipeNode.handleConnection() - Close connection with address %+v - Error: %v", addr, err)
		}
	}()
	data, err := ioutil.ReadAll(conn)
	if err == nil {
		pipe.outChan <- model.PipeMessage(data)
	} else {
		pipe.logger.Warnf("PipeNode.handleConnection() - Unread message from client %+v, Error %v", addr, err)
	}

}

func (pipe *pipeNode) callClient(message model.PipeMessage) {
	var err error
	var conn net.Conn
	pipe.registerClient()
	if pipe.config.Config == nil {
		// Plain connection
		conn, err = net.Dial("tcp", pipe.outputAddress)
	} else {
		// SSL/TLS Encryption
		conn, err = tls.Dial("tcp", pipe.outputAddress, pipe.config.Config)
	}
	if err != nil {
		pipe.logger.Errorf("PipeNode.callClient() - Error connecting with client %s: %v", pipe.outputAddress, err)
		return
	}
	defer func() {
		time.Sleep(1 * time.Second)
		err = conn.Close()
		if err != nil {
			pipe.logger.Errorf("PipeNode.callClient() - Error disconnecting from client %s: %v", pipe.outputAddress, err)
		}
		pipe.deregisterClient()
	}()
	_, err = conn.Write([]byte(message))
	if err != nil {
		pipe.logger.Errorf("PipeNode.callClient() - Error writing message to client %s: %v", pipe.outputAddress, err)
	}
	pipe.logger.Infof("PipeNode.callClient() - Message sent to client %s", pipe.outputAddress)
}

func (pipe *pipeNode) readFromInputChannel() {
	if ! pipe.running {
		pipe.running = true
	}
	pipe.inChan = make(chan model.PipeMessage)
	pipe.inChanCreated = true
	ClientCycle:
	for pipe.running {
		select {
		case msg := <- pipe.inChan:
			go pipe.callClient(msg)
		case <- time.After(ServerClientResetTimeout):
			if ! pipe.running {
				break ClientCycle
			}
			continue
		}
	}
}

func (pipe *pipeNode) acceptClients() {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("PipeNode.acceptClients() - Error: %v", r))
			pipe.logger.Fatalf("PipeNode.acceptClients() - Error: %v", err)
		}
	}()
	if pipe.tcpListener != nil {
		pipe.outChan = make(chan model.PipeMessage)
		pipe.outChanCreated = true
		for pipe.running{
			var conn net.Conn
			conn, err = (*pipe.tcpListener).Accept()
			if err != nil{
				pipe.logger.Errorf("PipeNode.acceptClients() - Acceptance Error: %v", err)
				continue
			}
			pipe.logger.Debugf("PipeNode.acceptClients() - Handling request from: %+v ...", conn.RemoteAddr())
			go pipe.handleConnection(conn)
		}
	} else {
		pipe.logger.Fatalf("PipeNode.acceptClients() - Invalid listener - Stopping server ...")
		err = pipe.Stop()
		pipe.logger.Error(err)
	}
}

func (pipe *pipeNode) shutdownTimer() {
	pipe.timer = time.NewTicker(5 * time.Second)
	defer func() {
		pipe.timer.Stop()
		pipe.timer = nil
	}()
	for {
		select {
		case <- pipe.timer.C:
			if pipe.checkExit() {
				pipe.evacuate()
				return
			}

		}
	}
}

func (pipe *pipeNode) registerRequest() {
	defer pipe.requestsMutex.Unlock()
	pipe.requestsMutex.Lock()
	pipe.activeRequests++
}

func (pipe *pipeNode) deregisterRequest() {
	defer pipe.requestsMutex.Unlock()
	pipe.requestsMutex.Lock()
	pipe.activeRequests--
}

func (pipe *pipeNode) registerClient() {
	defer pipe.clientsMutex.Unlock()
	pipe.clientsMutex.Lock()
	pipe.activeClients++
}

func (pipe *pipeNode) deregisterClient() {
	defer pipe.clientsMutex.Unlock()
	pipe.clientsMutex.Lock()
	pipe.activeClients--
}

func (pipe *pipeNode) working() bool {
	return pipe.activeRequests > 0 && pipe.activeClients > 0
}

func (pipe *pipeNode) checkExit() bool {
	if ! pipe.Running() && ! pipe.working() {
		pipe.internal <- exit
		return true
	}
	return false
}

func (pipe *pipeNode) Stop() error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("PipeNode.Stop() - Error: %v", r))
			pipe.logger.Fatalf("PipeNode.Stop() -  Error: %v", err)
		}
	}()
	if ! pipe.running {
		err = errors.New(fmt.Sprint("PipeNode.Stop() - Error: Server is already stopped"))
		pipe.logger.Errorf("PipeNode.Stop() -  %v", err)
		return err
	}
	pipe.internal <- shutdown
	go pipe.shutdownTimer()
	pipe.running = false
	if pipe.tcpListener != nil {
		err := (*pipe.tcpListener).Close()
		if err != nil {
			pipe.logger.Errorf("PipeNode.Stop() - Gently shutting down server error occurred: %v", err)
			pipe.logger.Warnf("PipeNode.Stop() - Try brute-force server close ...")
			pipe.activeRequests = 0
			pipe.activeClients = 0
		}
		pipe.tcpListener = nil
	}
	return err
}

func (pipe *pipeNode) Running() bool {
	return pipe.running
}

func (pipe *pipeNode) evacuate() {
	close(pipe.internal)
	pipe.internal = nil
	close(pipe.commands)
	pipe.commands = nil
	if pipe.config.Type == model.InputPipe || pipe.config.Type == model.InputOutputPipe {
		close(pipe.inChan)
		pipe.inChanCreated = false
		pipe.inChan = nil
	}
	if pipe.config.Type == model.OutputPipe || pipe.config.Type == model.InputOutputPipe {
		close(pipe.outChan)
		pipe.outChanCreated = false
		pipe.outChan = nil
	}
}

func (pipe *pipeNode) isOperating() bool {
	if pipe.config.Type == model.InputOutputPipe {
		return pipe.running && pipe.inChanCreated && pipe.outChanCreated
	} else if pipe.config.Type == model.InputPipe {
		return pipe.running && pipe.outChanCreated
	} else if pipe.config.Type == model.OutputPipe {
		return pipe.running && pipe.inChanCreated
	}
	return pipe.running
}

func (pipe *pipeNode) UntilStarted() {
	defer func() {
		if r := recover(); r != nil {
			pipe.logger.Fatalf("PipeNode.UntilStarted() - Fatal Error: %v", r)
		}
	}()
	pipe.logger.Debugf("PipeNode.UntilStarted() - Waiting for server running and input/output channel is open")
	for ! pipe.isOperating() {
		time.Sleep(ServerClientResetTimeout)
	}
	pipe.logger.Warnf("PipeNode.UntilStarted() - Server is now running and input/output pipe is open")
}

func (pipe *pipeNode) Wait() {
	defer func() {
		if r := recover(); r != nil {
			pipe.logger.Fatalf("PipeNode.Wait() - Fatal Error: %v", r)
		}
	}()
	pipe.logger.Debugf("PipeNode.Wait() - Waiting for server shutdown")
waitCycle:
	for pipe.Running() || pipe.working() {
		select {
		case sig := <- pipe.internal:
			if sig == shutdown {
				pipe.commands	<- purge
			} else if sig == exit {
				break waitCycle
			}
		case <- time.After(ServerWaitTimeout):
			continue
		}
	}
	pipe.logger.Warnf("PipeNode.Wait() - Server shutdown in progress, exiting ...")
	time.Sleep(10 * time.Second)
	pipe.logger.Warnf("PipeNode.Wait() - exit")
}

func (pipe *pipeNode) GetOutputPipeChannel() chan<- model.PipeMessage {
	return pipe.inChan
}

func (pipe *pipeNode) GetInputPipeChannel() <-chan model.PipeMessage {
	return pipe.outChan
}

func NewPipeNode(appName string, verbosity log.LogLevel) model.PipeNode {
	return &pipeNode{
		logger: log.NewLogger(appName, verbosity),
		requestsMutex: sync.Mutex{},
		clientsMutex: sync.Mutex{},
	}
}
