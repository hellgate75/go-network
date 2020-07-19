package model

import "crypto/tls"

type PipeType byte

const (
	// Input stream pipe
	NoTypeSelected		PipeType = 0
	// Input stream pipe
	InputPipe			PipeType = iota + 1
	// Output stream pipe
	OutputPipe
	// Input and Output streams pipe
	InputOutputPipe
)

// Default Message type
type PipeMessage []byte

// Describes an Pipe Node most features
type PipeNode interface {
	// Creates pipe node configuration, and setup the network properties.
	// It raises exception if the server is already running.
	Init(config PipeNodeConfig) (PipeNode, error)
	// Starts Pipe Node and serve requests
	Type() PipeType
	// Starts Pipe Node and serve requests
	Start() error
	// Stops Pipe Node and stop requests
	Stop() error
	// Verify Pipe Node is running
	Running() bool
	// Wait for Pipe Node is down
	Wait()
	// Wait for Pipe Node is started an output channel  (for Input or Input/Output Pipe mode nodes)
	UntilStarted()
	// Collects a message input channel (for Output or Input/Output Pipe mode nodes)
	GetOutputPipeChannel() chan<- PipeMessage
	// Collects a message output channel (for Input or Input/Output Pipe mode nodes)
	GetInputPipeChannel() <-chan PipeMessage
}

// Describe pine node properties
type PipeNodeConfig struct {
	// Input Host name or ip address (eg. my-host.acme.com or 127,0,0,1 or empty or 0.0.0.0)
	InHost 			string
	// Input Pipe Node Port
	InPort 			int
	// Output Host name or ip address (eg. my-host.acme.com or 127,0,0,1 or empty or 0.0.0.0)
	OutHost			string
	// Input Pipe Node Port
	OutPort 		int
	// Pipe Node type
	Type 			PipeType
	// Pipe Node Security Configuration
	Config 			*tls.Config
}
