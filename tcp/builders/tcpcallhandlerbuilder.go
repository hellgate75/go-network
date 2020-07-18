package builders

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model"
	context2 "github.com/hellgate75/go-network/model/context"
	"github.com/hellgate75/go-network/model/encoding"
	"github.com/hellgate75/go-network/tcp/stream"
	"net"
)

// Helper for creating a new model.TcpCallHandler
type TcpCallHandlerBuilder interface {
	// Define mandatory handler functions group name
	WithName(name string) TcpCallHandlerBuilder
	//Ass a new Method Handler to the builder
	WithTcpHandling(action model.TcpAction) TcpCallHandlerBuilder
	// Associate an error channel, for creating a flow of errors from the request
	WithErrorChannel(ch chan error) TcpCallHandlerBuilder
	// Build the model.TcpCallHandler and report any error occurred during the build process
	Build() (model.TcpCallHandler, error)
}

type tcpCallHandlerBuilder struct {
	name          string
	names         []string
	actions       []model.TcpAction
	errorHandling bool
	errCh         chan error
}

func (b *tcpCallHandlerBuilder) WithName(name string) TcpCallHandlerBuilder {
	b.name = name
	return b
}

func (b *tcpCallHandlerBuilder) WithTcpHandling(action model.TcpAction) TcpCallHandlerBuilder {
	b.actions = append(b.actions, action)
	b.names = append(b.names, action.GetName())
	return b
}
func (b *tcpCallHandlerBuilder) WithErrorChannel(ch chan error) TcpCallHandlerBuilder {
	b.errorHandling = true
	b.errCh = ch
	return b
}

func (b *tcpCallHandlerBuilder) Build() (model.TcpCallHandler, error) {
	var err error
	if len(b.name) == 0 {
		err = errors.New(fmt.Sprint("Empty name found"))
	}
	if len(b.actions) == 0 {
		err = errors.New(fmt.Sprint("No actions provided for the given name"))
	}
	return &tcpCallHandler{
		name:          b.name,
		names:         b.names,
		actions:       b.actions,
		errCh:         b.errCh,
		errorHandling: b.errorHandling,
		handlerMap:    make(map[string]interface{}),
	}, err
}


type  tcpCallHandler struct {
	names         []string
	name          string
	actions       []model.TcpAction
	errorHandling bool
	errCh         chan error
	handlerMap    map[string]interface{}
	serverMap     *map[string]interface{}
	encoding	  encoding.Encoding
	logger        log.Logger
}

func (h *tcpCallHandler) Names() []string {
	return h.names
}


func (h *tcpCallHandler) GetName() string {
	return h.name
}

func (h *tcpCallHandler) HandleRequest(conn net.Conn, closer stream.ConnReaderWriterCloser) {
	h.logger.Debugf("Running handler %s, waiting for data read ...", h.name)
	closer.Wait()
	h.logger.Debugf("Running handler %s, data has been read", h.name)
	for _, action := range h.actions {
		context := context2.NewTcpContext(conn, closer, h.encoding)
		// Set up reference to handler map cache element
		context.HandlerMap = &h.handlerMap
		context.Logger = h.logger
		// Set up reference to Tcp  global server map cache element
		context.ServerMap = h.serverMap
		err := action.With(context).Do()
		if err != nil && h.errorHandling {
			h.errCh <- err
		}
	}
}

func (h *tcpCallHandler) SetLogger(logger log.Logger) {
	h.logger = logger
}

func (h *tcpCallHandler) SetServerMap(m *map[string]interface{}) {
	h.serverMap = m
}
func (h *tcpCallHandler) SetEncoding(enc encoding.Encoding) {
	h.encoding = enc
}
func NewTcpCallHandlerBuilder() TcpCallHandlerBuilder {
	return &tcpCallHandlerBuilder{
		names: make([]string, 0),
		actions: make([]model.TcpAction, 0),
	}
}

