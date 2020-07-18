package model

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model/context"
	"net/http"
)

const (
	//Rest Web Method: GET
	REST_METHOD_GET RestMethod = http.MethodGet
	//Rest Web Method: POST
	REST_METHOD_POST RestMethod = http.MethodPost
	//Rest Web Method: POST (FORM)
	REST_METHOD_POST_FORM RestMethod = http.MethodPost + "_FORM"
	//Rest Web Method: HEAD
	REST_METHOD_HEAD RestMethod = http.MethodHead
	//Rest Web Method: CONNECT
	REST_METHOD_CONNECT RestMethod = http.MethodConnect
	//Rest Web Method: DELETE
	REST_METHOD_DELETE RestMethod = http.MethodDelete
	//Rest Web Method: OPTIONS
	REST_METHOD_OPTIONS RestMethod = http.MethodOptions
	//Rest Web Method: PATCH
	REST_METHOD_PATCH RestMethod = http.MethodPatch
	//Rest Web Method: PUT
	REST_METHOD_PUT RestMethod = http.MethodPut
	//Rest Web Method: TRACE
	REST_METHOD_TRACE RestMethod = http.MethodTrace
	
	//Rest Web Protocol: HTTP
	REST_PROTOCOL_HTTP RestProtocol = "http"
	//Rest Web Protocol: HTTPS
	REST_PROTOCOL_HTTPS RestProtocol = "https"
	//Rest Web Protocol: WS
	REST_PROTOCOL_WS RestProtocol = "ws"
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

// Rest Web Method  Type
type RestMethod string

// Rest Protocol Type
type RestProtocol string

// Defines an handler for an entire call
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

type ApiActionFunction func(context.ApiCallContext) error

