package api

import (
	"bytes"
	"errors"
	"fmt"
	io2 "github.com/hellgate75/go-network/io"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model"
	"github.com/hellgate75/go-network/model/encoding"
	"io"
	"io/ioutil"
	"net/http"
)

type apiClient struct{
	config 			*model.ClientConfig
	cli				*http.Client
	baseUrl			string
	logger			log.Logger
}

func (c *apiClient) Connect(config model.ClientConfig) error {
	c.config = &config
	if c.config.Protocol == "" || c.config.Host == "" || c.config.Port == 0 {
		c.logger.Fatal("Invalid protocol, server and/or port values")
		return errors.New(fmt.Sprint("Invalid protocol, server and/or port values"))
	}
	c.cli = &http.Client{
	}
	if c.config.Timeout > 0 {
		c.cli.Timeout = c.config.Timeout
	}
	c.baseUrl = fmt.Sprintf("%s://%s:%v", c.config.Protocol, c.config.Host, c.config.Port)
	c.logger.Debugf("Created default base url: %s", c.baseUrl)
	return nil
}


func (c *apiClient) Call(path string, method string, contentType *encoding.MimeType, accepts *encoding.MimeType, body io.Reader) (*http.Response, error) {
	if c.cli == nil {
		c.logger.Fatal("Client is not connected to a server socket")
		return nil, errors.New(fmt.Sprint("Client is not connected to a server socket"))
	}
	var err error
	var out *http.Response
	var r *http.Request
	var url = fmt.Sprintf("%s%s", c.baseUrl, path)
	c.logger.Debugf("Creating request with url: %s, and web method: %s ...", url, method)
	r, err = http.NewRequest(method, url, body)
	if err != nil {
		c.logger.Errorf("Error creating the request: %v", err)
		return out, err
	}
	if contentType != nil {
		c.logger.Debugf("Adding Content Type Header: %s", string(*contentType))
		r.Header.Add("Content-Type", string(*contentType))
	}
	if accepts != nil {
		c.logger.Debugf("Adding Accepts Header: %s", string(*accepts))
		r.Header.Add("Accepts", string(*accepts))
	}
	c.logger.Debug("Running the client handler using out request ...")
	return c.cli.Do(r)
}

func (c *apiClient) Encode(path string, method string, contentType encoding.MimeType, accepts *encoding.MimeType, request interface{}, response interface{}) error {
	if c.cli == nil {
		c.logger.Fatal("Client is not connected to a server socket")
		return errors.New(fmt.Sprint("Client is not connected to a server socket"))
	}
	var err error
	var r *http.Request
	var data []byte
	var url = fmt.Sprintf("%s%s", c.baseUrl, path)
	c.logger.Debugf("Creating request with url: %s, and web method: %s ...", url, method)
	requestEncoding := encoding.ParseMimeType(contentType)
	data, err = io2.Marshal(requestEncoding, request)
	if err != nil {
		c.logger.Errorf("Error parsing the request mime type to encoding.Encoding: %v", err)
		return err
	}
	r, err = http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		c.logger.Errorf("Error creating the request: %v", err)
		return err
	}
	if string(contentType) != "" {
		c.logger.Debugf("Adding Content Type Header: %s", string(contentType))
		r.Header.Add("Content-Type", string(contentType))
	}
	if accepts != nil {
		c.logger.Debugf("Adding Accepts Header: %s", string(*accepts))
		r.Header.Add("Accepts", string(*accepts))
	}
	c.logger.Debug("Running the client handler using out request ...")
	resp, err :=  c.cli.Do(r)
	if err != nil {
		c.logger.Errorf("Error sending the request: %v", err)
		return errors.New(fmt.Sprintf("Error sending the request: %v", err))
	}
	if accepts != nil && response != nil {
		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.logger.Errorf("Error reading the response body: %v", err)
			return errors.New(fmt.Sprintf("Error reading the response body: %v", err))
		}
		responseEncoding := encoding.ParseMimeType(*accepts)
		err = io2.Unmarshal(respData, responseEncoding, response)
	}
	return err
}

func NewApiClient(appName string, verbosity log.LogLevel) model.ApiClient {
	return &apiClient{
		logger: log.NewLogger(appName, verbosity),
	}
}
