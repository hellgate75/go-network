package api

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-network/model"
	"github.com/hellgate75/go-network/model/encoding"
	"io"
	"net/http"
)

type apiClient struct{
	config 			*model.ClientConfig
	cli				*http.Client
	baseUrl			string
}

func (c *apiClient) Connect(config model.ClientConfig) error {
	c.config = &config
	if c.config.Protocol == "" || c.config.Host == "" || c.config.Port == 0 {
		return errors.New(fmt.Sprint("Invalid protocol, server and/or port values"))
	}
	c.cli = &http.Client{
	}
	if c.config.Timeout > 0 {
		c.cli.Timeout = c.config.Timeout
	}
	c.baseUrl = fmt.Sprintf("%s://%s:%v", c.config.Protocol, c.config.Host, c.config.Port)
	return nil
}

func (c *apiClient) Call(path string, method string, contentType *encoding.MimeType, accepts *encoding.MimeType, body io.Reader) (*http.Response, error) {
	var err error
	var out *http.Response
	var r *http.Request
	r, err = http.NewRequest(method, fmt.Sprintf("%s%s", c.baseUrl, path), body)
	if err != nil {
		return out, err
	}
	if contentType != nil {
		r.Header.Add("Content-Type", string(*contentType))
	}
	if accepts != nil {
		r.Header.Add("Accepts", string(*accepts))
	}
	return c.cli.Do(r)
}

func NewApiClient() model.ApiClient {
	return &apiClient{
	}
}