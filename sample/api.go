package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hellgate75/go-network/api"
	"github.com/hellgate75/go-network/api/builders"
	"github.com/hellgate75/go-network/io"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model/context"
	"github.com/hellgate75/go-network/model/encoding"
	"io/ioutil"
	"net/http"
)

func sampleStruct() interface{} {
	return struct{
		Id			string		`yaml:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
		Name		string		`yaml:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
		Surname		string		`yaml:"surname,omitempty" json:"surname,omitempty" xml:"surname,omitempty"`
		Age			int			`yaml:"age,omitempty" json:"age,omitempty" xml:"age,omitempty"`
	}{
		"1",
		"Fabrizio",
		"Torelli",
		45,
	}
}

func emptyStruct() interface{} {
	return struct{
		Id			string		`yaml:"id,omitempty" json:"id,omitempty" xml:"id,omitempty"`
		Name		string		`yaml:"name,omitempty" json:"name,omitempty" xml:"name,omitempty"`
		Surname		string		`yaml:"surname,omitempty" json:"surname,omitempty" xml:"surname,omitempty"`
		Age			int			`yaml:"age,omitempty" json:"age,omitempty" xml:"age,omitempty"`
	}{
	}
}


func TestApiServer() {
	apiServer := api.NewApiServer("Sample Api Rest Server", log.DEBUG)
	serverConfig, err := builders.
		NewServerConfigBuilder().
		WithHost("", 9999).
		Build()
	if err != nil {
		panic(err)
	}
	apiServer, err = apiServer.Init(serverConfig)
	if err != nil {
		panic(err)
	}
	handler, _ := builders.NewApiCallHandlerBuilder().
		WithPath("/").
		WithWebMethodHandling("POST", builders.NewApiActionBuilder().
			With(func(c context.ApiCallContext) error {
				defer func() {
					if r := recover(); r != nil {
						err = errors.New(fmt.Sprintf("%v", r))
					}
				}()
				req := emptyStruct()
				err = c.ParseBody(&req)
				c.Logger.Infof("POST Data Read Error: %v", err)
				c.Logger.Infof("POST Data: %+v", req)
				res := sampleStruct()
				err = c.WriteResponse(&res, http.StatusOK)
				return err
			}).
			Build()).
		Build()
	err = apiServer.AddPath(handler)
	if err != nil {
		panic(err)
	}
	err = apiServer.Start()
	if err != nil {
		panic(err)
	}
	apiServer.Wait()
}

func TestApiClient() {
	logger := log.NewLogger("Sample Api Rest Server", log.DEBUG)
	apiClient := api.NewApiClient("Sample Api Rest Server", log.DEBUG)
	apiClientConfig, err := builders.NewClientConfigBuilder().
		WithHost("http", "localhost", 9999).
		Build()
	if err != nil {
		panic(err)
	}
	err = apiClient.Connect(apiClientConfig)
	if err != nil {
		panic(err)
	}
	empty := emptyStruct()
	b, err := io.Marshal(encoding.EncodingJSONFormat, sampleStruct())
	if err != nil {
		panic(err)
	}
	reader := bytes.NewBuffer(b)
	mime := encoding.JsonMimeType
	logger.Infof("Request: %s", string(b))
	resp, err := apiClient.Call("/", http.MethodPost, &mime, &mime, reader)
	logger.Infof("Status: %s", resp.Status)
	bts, _ := ioutil.ReadAll(resp.Body)
	logger.Infof("Response: %+v", resp)
	_ = io.Unmarshal(bts, encoding.EncodingJSONFormat, &empty)
	logger.Infof("Response data: %+v", empty)
}
