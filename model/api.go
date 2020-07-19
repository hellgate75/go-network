package model

import (
	"crypto/tls"
	"github.com/hellgate75/go-network/model/encoding"
	"io"
	"net/http"
	"time"
)

// Describes an API Server most features
type ApiServer interface {
	// Create server configuration, and setup the network properties.
	// It raises exception if the server is already running.
	Init(config ServerConfig) (ApiServer, error)
	// Start API Server and serve requests
	Start() error
	// Stop API Server and stop requests
	Stop() error
	// Verify API Server is running
	Running() bool
	// Verify API Server is running
	Working() bool
	// Wait for API server is down
	Wait()
	// Add a new path call handler in the api router, allowing management of multiple
	// mime types and methods calls for the same path requested by the client
	// It raises exception if the API call handler has not method call handling function
	// or if the Path is duplicate
	AddPath(ApiCallHandler) error
}

// Describes an API Client most features
type ApiClient interface {
	// Configure a new Connection using server base path
	Connect(config ClientConfig) error
	// Make a call
	// Request must be sent to the body Reader (preferred: bytes.Buffer)
	Call(path string, method string, contentType *encoding.MimeType, accepts *encoding.MimeType, body io.Reader) (*http.Response, error)
	// Make a call
	// Request must be sent and object with preferred encoding configuration
	Encode(path string, method string, contentType encoding.MimeType, accepts *encoding.MimeType, request interface{}, response interface{}) error
}

// Describe client connection properties
type ClientConfig struct {
	// Communication protocol (eg.: http, https, ...)
	Protocol 	string
	// Host name or ip address (eg. my-host.acme.com or 192.168.1.222)
	Host 		string
	// Remote API Server Port
	Port 		int
	// Remote API Server connection timeout (0 means not set)
	Timeout		time.Duration
	// Remote API Server Security Configuration
	Config 		*tls.Config
}


// Describe server connection properties
type ServerConfig struct {
	// Host name or ip address (eg. my-host.acme.com or 127,0,0,1 or empty or 0.0.0.0)
	Host 		string
	// API Server Port
	Port 		int
	// API Server Security Configuration
	Config 		*tls.Config
	// TLS Certificate file Full Path
	CertPath	string
	// TLS Certificate Key file Full Path
	KeyPath		string
}
