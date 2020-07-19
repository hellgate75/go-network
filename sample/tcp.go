package main

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model/context"
	"github.com/hellgate75/go-network/model/encoding"
	"github.com/hellgate75/go-network/tcp"
	"github.com/hellgate75/go-network/tcp/builders"
	"time"
)


func TestTcpServer() {
	logger := log.NewLogger("Sample Tcp Server", log.DEBUG)
	tcpServer := tcp.NewTcpServer("Sample Tcp Server", log.DEBUG)
	serverConfig, err := builders.
		NewTcpServerConfigBuilder().
		WithNetwork("tcp").
		WithEncoding(encoding.EncodingJSONFormat).
		WithHost("", 9998).
		Build()
	if err != nil {
		panic(err)
	}
	logger.Infof("Server Config Encoding: %v", serverConfig.Encoding)
	tcpServer, err = tcpServer.Init(serverConfig)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	err = tcpServer.Start()
	if err != nil {
		panic(err)
	}
	tcpServer.Wait()
}

func TestTcpClient() {
	logger := log.NewLogger("Sample Tcp Client", log.DEBUG)
	tcpClient := tcp.NewTcpClient("Sample Tcp Client", log.DEBUG)
	tcpClientConfig, err := builders.NewTcpClientConfigBuilder().
		WithHost("localhost", 9998).
		WithNetwork("tcp").
		WithEncoding(encoding.EncodingJSONFormat).
		Build()
	if err != nil {
		panic(err)
	}
	logger.Info("Connecting server ...")
	err = tcpClient.Connect(tcpClientConfig)
	if err != nil {
		panic(err)
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
