package context

import (
	io2 "github.com/hellgate75/go-network/io"
	"github.com/hellgate75/go-network/log"
	"github.com/hellgate75/go-network/model/encoding"
	"io"
	"io/ioutil"
	"net"
)

// Describe Tcp transaction context
type TcpContext struct {
	// Unique request identifier
	Id string
	// Connection Response Writer component
	ResponseWriter io.Writer
	// Connection Request reader component
	RequestReader io.Reader
	// Remote host address
	RemoteAddress net.Addr
	// Request MimeType
	ResponseEncoding encoding.Encoding
	// Response MimeType
	RequestEncoding encoding.Encoding
	// Request level cache map element
	RequestMap map[string]interface{}
	// Reference to Handler level cache map element
	HandlerMap *map[string]interface{}
	// Reference to Api Server level cache map element
	ServerMap *map[string]interface{}
	// Reference to Api Server level cache map element
	Logger log.Logger
}

func NewTcpContext(conn net.Conn, reader io.Reader, serverEncoding encoding.Encoding) TcpContext {
	return TcpContext{
		Id:               GenerateUUUID(),
		ResponseWriter:	  conn,
		RequestReader:	  reader,
		RemoteAddress:    conn.RemoteAddr(),
		RequestEncoding:  serverEncoding,
		ResponseEncoding:  serverEncoding,
		RequestMap:       make(map[string]interface{}),
		HandlerMap:       nil,
		ServerMap:        nil,

	}
}

func (ctx *TcpContext) ParseRequest(requestBody interface{}) error {
	data, err := ioutil.ReadAll(ctx.RequestReader)
	if err != nil {
		return err
	}
	err = io2.Unmarshal(data, ctx.RequestEncoding, requestBody)
	return err
}


func (ctx *TcpContext) WriteResponse(responseBody interface{}) error {
	data, err := io2.Marshal(ctx.ResponseEncoding, responseBody)
	if err != nil {
		return err
	}
	_, err = ctx.ResponseWriter.Write(data)
	return err
}
