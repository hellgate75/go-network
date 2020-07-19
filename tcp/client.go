package tcp

import (
	"crypto/tls"
	"errors"
	"fmt"
	io2 "github.com/hellgate75/go-network/io"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model"
	"io"
	"io/ioutil"
	"net"
	"time"
)

type tcpClient struct{
	config 			*model.TcpClientConfig
	cli				net.Conn
	logger			log.Logger
}

func (c *tcpClient) Connect(config model.TcpClientConfig) error {
	var err error
	var conn net.Conn
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("TcpClient.Connect() - Error: %v", r))
			c.logger.Fatal(err)
		}
	}()
	if c.IsOpen() {
		err := c.Close()
		if err != nil {
			return err
		}

	}
	c.config = &config
	if c.config.Network == "" {
		c.logger.Error("Invalid network value")
		return errors.New(fmt.Sprint("Invalid network, server and/or port values"))
	}
	address := fmt.Sprintf("%s:%v", c.config.Host, c.config.Port)
	if c.config.Port <= 0 {
		address = fmt.Sprintf("%s", c.config.Host)
	}

	if config.Config == nil {
		// Plain connection
		conn, err = net.Dial(config.Network, address)
	} else {
		// SSL/TLS Encryption
		conn, err = tls.Dial(config.Network, address, config.Config)
	}
	if err == nil {
		if c.config.Timeout > 0 {
			err = conn.SetDeadline(time.Now().Add(c.config.Timeout))
		}
		c.cli = conn
	} else {
		c.logger.Error(err)
	}
	return err
}

func (c *tcpClient) IsOpen() bool {
	return c.cli != nil
}

func (c *tcpClient) Close() error {
	if c.cli == nil {
		c.logger.Error("Connection is already closed ...")
		return errors.New(fmt.Sprint("Connection is already closed ..."))
	}
	return c.cli.Close()
}
func (c *tcpClient) readParseInput(response interface{}) error {
	if response != nil {
		data, err := ioutil.ReadAll(c.cli)
		if err != nil {
			return err
		}
		err = io2.Unmarshal(data, c.config.Encoding, response)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *tcpClient) Send(body io.Reader, response interface{}, timeout time.Duration) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("TcpClient.Send() - Error: %v", r))
			c.logger.Fatal(err)
		}
	}()
	if c.cli == nil {
		c.logger.Error("Client is not connected to a server socket")
		return errors.New(fmt.Sprint("Client is not connected to a server socket"))
	}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	c.logger.Debug("Sending data to client ...")
	_, err = c.cli.Write(data)
	if err != nil {
		return err
	}
	c.logger.Debug("Waiting server received the request...")
	time.Sleep(timeout)
	if response != nil {
		c.logger.Debug("Reading for answer...")
		err =  c.readParseInput(response)
		if err != nil {
			c.logger.Error(err)
		}
	}
	return err
}

func (c *tcpClient) Encode(request interface{}, response interface{}, timeout time.Duration) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("TcpClient.Encode() - Error: %v", r))
			c.logger.Fatal(err)
		}
	}()
	if c.cli == nil {
		c.logger.Error("Client is not connected to a server socket")
		return errors.New(fmt.Sprint("Client is not connected to a server socket"))
	}
	var data []byte
	data, err = io2.Marshal(c.config.Encoding, request)
	if err != nil {
		return err
	}
	_, err = c.cli.Write(data)
	if err != nil {
		return err
	}
	c.logger.Debug("Waiting server received the request...")
	time.Sleep(timeout)
	if response != nil {
		c.logger.Debug("Reading for answer...")
		err =  c.readParseInput(response)
		if err != nil {
			c.logger.Error(err)
		}
	}
	return err
}
func (c *tcpClient) ReadRemote(timeout time.Duration, response interface{}) error {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("TcpClient.ReadRemote() - Error: %v", r))
			c.logger.Fatal(err)
		}
	}()
	if c.cli == nil {
		c.logger.Error("Client is not connected to a server socket")
		return errors.New(fmt.Sprint("Client is not connected to a server socket"))
	}
	if response != nil {
		var start = time.Now()
	readCycle:
		for timeout == 0 || time.Now().Sub(start) < timeout {
			err =  c.readParseInput(response)
			if err != nil {
				break readCycle
			}
			time.Sleep(2 * time.Second)
		}
	} else {
		err = errors.New(fmt.Sprint("Nil response interface, cannot parse the remote connection stream"))
	}
	return err
}

func NewTcpClient(appName string, verbosity log.LogLevel) model.TcpClient {
	return &tcpClient{
		logger: log.NewLogger(appName, verbosity),
	}
}