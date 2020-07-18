package builders

import (
	"errors"
	"fmt"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model"
	context2 "github.com/hellgate75/go-network/model/context"
	"net/http"
	"strings"
)

// Helper for creating a new model.ApiCallHandler
type ApiCallHandlerBuilder interface {
	// Define mandatory path in various formats (eg.: /path/ or /path/{var} etc...)
	// accordingly to the gorilla mux Router specifications
	WithPath(path string) ApiCallHandlerBuilder
	//Ass a new Method Handler to the builder
	WithWebMethodHandling(method string, action model.ApiAction) ApiCallHandlerBuilder
	// Associate an error channel, for creating a flow of errors from the request
	WithErrorChannel(ch chan error) ApiCallHandlerBuilder
	// Build the model.ApiCallHandler and report any error occurred during the build process
	Build() (model.ApiCallHandler, error)
}

type apiCallHandlerBuilder struct {
	path		string
	methods		map[string]model.ApiAction
	errorHandling	bool
	errCh			chan error
}

func (b *apiCallHandlerBuilder) WithPath(path string) ApiCallHandlerBuilder {
	b.path = path
	return b
}

func (b *apiCallHandlerBuilder) WithWebMethodHandling(method string, action model.ApiAction) ApiCallHandlerBuilder {
	if method!= "" {
		var m = strings.ToUpper(method)
		b.methods[m] = action
	}
	return b
}
func (b *apiCallHandlerBuilder) WithErrorChannel(ch chan error) ApiCallHandlerBuilder {
	b.errorHandling = true
	b.errCh = ch
	return b
}

func (b *apiCallHandlerBuilder) Build() (model.ApiCallHandler, error) {
	var err error
	var methods = make([]string, 0)
	for k, _ := range b.methods {
		methods = append(methods, k)
	}
	if len(b.path) == 0 {
		err = errors.New(fmt.Sprint("Empty path found"))
	}
	if len(methods) == 0 {
		err = errors.New(fmt.Sprint("No methods provided for the given path"))
	}
	return &apiCallHandler{
		path: b.path,
		methods: methods,
		actions: b.methods,
		errCh: b.errCh,
		errorHandling: b.errorHandling,
		handlerMap: make(map[string]interface{}),
	}, err
}


type  apiCallHandler struct {
	methods 		[]string
	path			string
	actions			map[string]model.ApiAction
	errorHandling	bool
	errCh			chan error
	handlerMap 		map[string]interface{}
	serverMap 		*map[string]interface{}
	logger			log.Logger
}

func (h *apiCallHandler) Methods() []string {
	return h.methods
}

func (h *apiCallHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	m := strings.ToUpper(r.Method)
	var action model.ApiAction
	var ok bool
	if action, ok = h.actions[m]; !ok {
		http.NotFound(w, r)
	} else {
		context := context2.NewApiCallContext(w, r)
		// Set up reference to handler map cache element
		context.HandlerMap = &h.handlerMap
		context.Logger = h.logger
		// Set up reference to Api  global server map cache element
		context.ServerMap = h.serverMap
		err := action.With(context).Do()
		if err != nil && h.errorHandling {
			h.errCh <- err
		}
	}
}

func (h *apiCallHandler) GetPath() string {
	return h.path
}

func (h *apiCallHandler) SetLogger(logger log.Logger) {
	h.logger = logger
}

func (h *apiCallHandler) SetServerMap(m *map[string]interface{}) {
	h.serverMap = m
}

func NewApiCallHandlerBuilder() ApiCallHandlerBuilder {
	return &apiCallHandlerBuilder{
		methods: make(map[string]model.ApiAction),
	}
}

