package main

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model/context"
	"github.com/hellgate75/go-network/model/encoding"
	"github.com/hellgate75/go-network/tcp"
	"github.com/hellgate75/go-network/tcp/builders"
	log2 "log"
	"time"
)


func TestTcpServer() {
	logger := log.NewLogger("Sample Tcp", log.DEBUG)
	tcpServer := tcp.NewTcpServer("Sample TCP", log.DEBUG)
	serverConfig, err := builders.
		NewTcpServerConfigBuilder().
		WithNetwork("tcp").
		WithEncoding(encoding.EncodingJSONFormat).
		WithHost("", 9998).
		Build()
	if err != nil {
		log2.Fatal(err)
	}
	logger.Infof("Server Config Encoding: %v", serverConfig.Encoding)
	tcpServer, err = tcpServer.Init(serverConfig)
	if err != nil {
		log2.Fatal(err)
	}
	handler, _ := builders.NewTcpCallHandlerBuilder().
		WithName("Main").
		WithTcpHandling(builders.NewTcpActionBuilder().
			WithName("read-sample-data").
			With(func(c context.TcpContext) error {
				logger.Info("Request handler Main.read-sample-data")
				defer func() {
					if r := recover(); r != nil {
						err = errors.New(fmt.Sprintf("%v", r))
					}
				}()
				logger.Debug("Reading request ...")
				req := emptyStruct()
				err = c.ParseRequest(&req)
				c.Logger.Debugf("Main Data Read Error: %v", err)
				c.Logger.Debugf("Main Data: %+v", req)
				res := sampleStruct()
				c.Logger.Debug("Writing the answer ...")
				err = c.WriteResponse(&res)
				if err != nil {
					time.Sleep(1 * time.Second)
					c.Logger.Errorf("Error during answer sending procedure: %v", err)
				} else {
					c.Logger.Info("Answer sent!!")
				}
				return err
			}).
			Build()).
		Build()
	err = tcpServer.AddPath(handler)
	if err != nil {
		log2.Fatal(err)
	}
	err = tcpServer.Start()
	if err != nil {
		log2.Fatal(err)
	}
	tcpServer.Wait()
}

func TestTcpClient() {
	logger := log.NewLogger("Client Tcp", log.DEBUG)
	tcpClient := tcp.NewTcpClient("Client Tcp", log.DEBUG)
	tcpClientConfig, err := builders.NewTcpClientConfigBuilder().
		WithHost("localhost", 9998).
		WithNetwork("tcp").
		WithEncoding(encoding.EncodingJSONFormat).
		Build()
	if err != nil {
		log2.Fatal(err)
	}
	logger.Info("Connecting server ...")
	err = tcpClient.Connect(tcpClientConfig)
	if err != nil {
		log2.Fatal(err)
	}
	defer func() {
		logger.Info("Closing connection ...")
		_ = tcpClient.Close()
	}()
	empty := emptyStruct()
	sample := sampleStruct()
	time.Sleep(1 * time.Second)
	logger.Infof("Request data: %v", sample)
	err = tcpClient.Encode(&sample, &empty, 0 * time.Second)
	logger.Infof("Response data: %+v", empty)
}
