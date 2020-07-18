package context

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hellgate75/go-network/io"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model/encoding"
	"io/ioutil"
	"net/http"
	"strings"
)

// Defines the API Call Context, used as ApiAction single source of  truth information
type ApiCallContext struct {
	// Unique request identifier
	Id string
	// Request path
	Path string
	// Request web method
	Method string
	// Connection Response Writer component
	ResponseWriter http.ResponseWriter
	// Connection Request component
	Request *http.Request
	// Request MimeType
	ContentMimeType encoding.MimeType
	// Response MimeType
	ResponseMimeType encoding.MimeType
	// Request level cache map element
	RequestMap map[string]interface{}
	// Reference to Handler level cache map element
	HandlerMap *map[string]interface{}
	// Reference to Api Server level cache map element
	ServerMap *map[string]interface{}
	// Reference to Api Server level cache map element
	Logger log.Logger
}

func NewApiCallContext(w http.ResponseWriter,
	r *http.Request) ApiCallContext {
	return ApiCallContext{
		Id:               GenerateUUUID(),
		Path:             r.URL.Path,
		Method:           strings.ToUpper(r.Method),
		ResponseWriter:   w,
		Request:          r,
		ContentMimeType:  getContentMime(r.Header),
		ResponseMimeType: getResponseMime(r.Header),
		RequestMap:       make(map[string]interface{}),
		HandlerMap:       nil,
		ServerMap:        nil,
	}
}

func getContentMime(header http.Header) encoding.MimeType {
	tp := header.Get("content-type")
	if tp == "" {
		tp = header.Get("Content-Type")
	}
	if tp == "" {
		return encoding.JsonMimeType
	}
	return encoding.MimeType(tp)
}

func getResponseMime(header http.Header) encoding.MimeType {
	tp := header.Get("accepts")
	if tp == "" {
		tp = header.Get("Accepts")
	}
	if tp == "" {
		return encoding.JsonMimeType
	}
	return encoding.MimeType(tp)
}

func (ctx *ApiCallContext) RequestEncoding() encoding.Encoding {
	return encoding.ParseMimeType(ctx.ContentMimeType)
}

func (ctx *ApiCallContext) ResponseEncoding() encoding.Encoding {
	return encoding.ParseMimeType(ctx.ResponseMimeType)
}

func (ctx *ApiCallContext) CanParseBody() bool {
	return ctx.Method == "POST"
}

func (ctx *ApiCallContext) ParseBody(requestBody interface{}) error {
	if ctx.Method != "POST" {
		return errors.New(fmt.Sprintf("Invalid web method: %s for requesting body parsing", ctx.Method))
	}
	var encodingValue = encoding.ParseMimeType(ctx.ContentMimeType)
	if ctx.RequestEncoding() == encoding.EncodingUNKNOWNFormat {
		return errors.New(fmt.Sprintf("Unable to discover an encoder for mime type: %v", ctx.ContentMimeType))
	}
	data, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	err = io.Unmarshal(data, encodingValue, requestBody)
	return err
}

func (ctx *ApiCallContext) WriteResponse(responseBody interface{}, code int) error {
	if ctx.ResponseEncoding() == encoding.EncodingUNKNOWNFormat {
		return errors.New(fmt.Sprintf("Unable to discover an encoder for mime type: %v", ctx.ResponseMimeType))
	}
	data, err := io.Marshal(ctx.ResponseEncoding(), responseBody)
	if err != nil {
		return err
	}
	ctx.ResponseWriter.WriteHeader(code)
	_, err = ctx.ResponseWriter.Write(data)
	return err
}

// Generate a Security Token of a given length
func GenerateUUUID() string {
	return uuid.New().String()
}

