package main

import (
	"fmt"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model"
	"github.com/hellgate75/go-network/pipe"
	"github.com/hellgate75/go-network/pipe/builders"
	"time"
)


func TestInputPipeNode() {
	logger := log.NewLogger("Sample Input Pipe Node", log.DEBUG)
	pipeInputNode := pipe.NewPipeNode("Sample Input Pipe Node", log.DEBUG)
	pipeNodeConfig, err := builders.
		NewPipeNodeConfigBuilder().
		WithInHost("", 9997).
		Build()
	if err != nil {
		panic(err)
	}
	pipeInputNode, err = pipeInputNode.Init(pipeNodeConfig)
	if err != nil {
		panic(err)
	}
	err = pipeInputNode.Start()
	if err != nil {
		panic(err)
	}
	defer func() {
		err = pipeInputNode.Stop()
		if err != nil {
			panic(err)
		}
	}()
	pipeInputNode.UntilStarted()
	outputMessageChannel := pipeInputNode.GetInputPipeChannel()
	logger.Debug("Start message network reader ...")
	for {
		select {
			case msg := <-outputMessageChannel:
				logger.Warnf("Received message: %s", string([]byte(msg)))
		}
	}
}

func TestOutputPipeNode() {
	logger := log.NewLogger("Sample Output Pipe Node", log.DEBUG)
	pipeOutputNode := pipe.NewPipeNode("Sample Output Pipe Node", log.DEBUG)
	pipeNodeConfig, err := builders.
		NewPipeNodeConfigBuilder().
		WithOutHost("", 9997).
		Build()
	if err != nil {
		panic(err)
	}
	pipeOutputNode, err = pipeOutputNode.Init(pipeNodeConfig)
	if err != nil {
		panic(err)
	}
	err = pipeOutputNode.Start()
	if err != nil {
		panic(err)
	}
	defer func() {
		err = pipeOutputNode.Stop()
		if err != nil {
			panic(err)
		}
	}()
	pipeOutputNode.UntilStarted()
	inputMessageChannel := pipeOutputNode.GetOutputPipeChannel()
	for {
		msg := createPipeMessage()
		logger.Debugf("Sending message: %s", string([]byte(msg)))
		inputMessageChannel <- msg
		logger.Warnf("Sent message: %s", string([]byte(msg)))
		time.Sleep(2 * time.Second)
	}
}

var pipeMessageCount int64

func createPipeMessage() model.PipeMessage {
	pipeMessageCount++
	return model.PipeMessage([]byte(fmt.Sprintf("This is message # %v", pipeMessageCount)))
}