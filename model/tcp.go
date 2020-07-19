package model

import (
	"crypto/tls"
	"github.com/hellgate75/go-network/model/encoding"
	"io"
	"time"
)

// Describes an Tcp Server most features
type TcpServer interface {
	// Creates server configuration, and setup the network properties.
	// It raises exception if the server is already running.
	Init(config TcpServerConfig) (TcpServer, error)
	// Starts Tcp Server and serve requests
	Start() error
	// Stops Tcp Server and stop requests
	Stop() error
	// Verifies Tcp Server is running
	Running() bool
	// Verifies Tcp Server is running
	Working() bool
	// Waits for Tcp server is down
	Wait()
	// Add a new path call handler in the api router, allowing management of multiple
	// mime types and methods calls for the same path requested by the client
	// It raises exception if the API call handler has not method call handling function
	// or if the Path is duplicate
	AddPath(TcpCallHandler) error
}

// Describes an Tcp Client most features
type TcpClient interface {
	// Configure a new Connection using server base path
	Connect(config TcpClientConfig) error
	// Close client connection
	Close() error
	// check if client connection is open
	IsOpen() bool
	// Make a call
	// Request must be sent to the body Reader (preferred: bytes.Buffer)
	Send(body io.Reader, response interface{}, timeout time.Duration) error
	// Make a call
	// Request must be sent and object with preferred encoding configuration
	Encode(request interface{}, response interface{}, timeout time.Duration) error
	// Wait for a client answer, for the maximum timeout of forever in case the timeout is zero
	ReadRemote(timeout time.Duration, response interface{}) error
}

// Describe client connection properties
type TcpClientConfig struct {
	// Connection network type (default: tcp)
	Network		string
	// Host name or ip address (eg. my-host.acme.com or 192.168.1.222)
	Host 		string
	// Remote Tcp Server Port
	Port 		int
	// Remote Tcp Server connection timeout (0 means not set)
	Timeout		time.Duration
	// Remote Tcp Server Security Configuration
	Config 		*tls.Config
	// Encoding
	Encoding		encoding.Encoding
}


// Describe server connection properties
type TcpServerConfig struct {
	// Connection network type (default: tcp)
	Network		string
	// Host name or ip address (eg. my-host.acme.com or 127,0,0,1 or empty or 0.0.0.0)
	Host 		string
	// Tcp Server Port
	Port 		int
	// Tcp Server Security Configuration
	Config 		*tls.Config
	// Encoding
	Encoding		encoding.Encoding
}
