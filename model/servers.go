package model

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model/context"
	"github.com/hellgate75/go-network/model/encoding"
	"github.com/hellgate75/go-network/tcp/stream"
	"net"
	"net/http"
)

// Context Key Type
type ContextKey string

func (c ContextKey) String() string {
	return "context key " + string(c)
}

var (
	// Session Context Key
	ContextSessionKey = ContextKey("session-key")
	// Session Context Auth Token
	ContextKeyAuthtoken = ContextKey("auth-token")
	// Session Context Remote Address
	ContextRemoteAddress = ContextKey("remote-address")
)

// Generate a Security Token of a given length
func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// Submit Success Data To the client
func SubmitSuccess(w http.ResponseWriter, message string) {
	w.WriteHeader(200)
	_, _ = w.Write([]byte(message))
}

// Submit Failure Data To the client
func SubmitFaiure(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_, _ = w.Write([]byte(message))
}

// Defines an handler for all methods in a request
type ApiCallHandler interface {
	// Returns the list of managed Web Methods [GET, POST, PUT, DELETE, ...]
	Methods() []string
	// Handle Request using http.Request, http.ResponseWriter
	HandleRequest(http.ResponseWriter, *http.Request)
	// Returns the path filter for this https call handler
	GetPath() string
	// Set reference to server map or leave map nil, if not used
	SetServerMap(m *map[string]interface{})
	// Set the server logger
	SetLogger(logger log.Logger)
}

// Interface that describes the callback action of an API call
type ApiAction interface {
	// Execute API command with API given arguments
	With(context.ApiCallContext) ApiAction
	Do() error
}

// Describe execution function for API rest connections
type ApiActionFunction func(context.ApiCallContext) error


// Defines an handler for an multiple actionsin a request
type TcpCallHandler interface {
	// Returns the list of managed actions names
	Names() []string
	// Handle Request using net.Conn
	HandleRequest(net.Conn, stream.ConnReaderWriterCloser)
	// Returns the handler name
	GetName() string
	// Set reference to server map or leave map nil, if not used
	SetServerMap(m *map[string]interface{})
	// Set the server logger
	SetLogger(logger log.Logger)
	// Set encoding used by the server
	SetEncoding(enc encoding.Encoding)
}

// Interface that describes the callback action of an Tcp request
type TcpAction interface {
	GetName() string
	// Execute Tcp command with API given arguments
	With(context.TcpContext) TcpAction
	Do() error
}


// Describe execution function for tcp connections
type TcpActionFunction func(context.TcpContext) error
