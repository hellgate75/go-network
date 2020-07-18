package builders

import (
	"github.com/hellgate75/go-network/model"
	"github.com/hellgate75/go-network/model/context"
)

type TcpActionBuilder interface {
	WithName(name string) TcpActionBuilder
	With(function model.TcpActionFunction) TcpActionBuilder
	Build() model.TcpAction
}

type tcpActionBuilder struct {
	function			model.TcpActionFunction
	name				string
}

func (builder *tcpActionBuilder) WithName(name string) TcpActionBuilder {
	builder.name = name
	return builder
}

func (builder *tcpActionBuilder) With(function model.TcpActionFunction) TcpActionBuilder {
	builder.function = function
	return builder
}

func (builder *tcpActionBuilder) Build() model.TcpAction {
	return &tcpAction{
		function: builder.function,
		name: builder.name,
	}
}

type tcpAction struct{
	function model.TcpActionFunction
	context *context.TcpContext
	name	string
}
func (action *tcpAction) GetName() string {
	return action.name
}
func (action *tcpAction) With(context context.TcpContext) model.TcpAction {
	action.context = &context
	return action
}

func (action *tcpAction) Do() error {
	return action.function(*action.context)
}

func NewTcpActionBuilder() TcpActionBuilder {
	return &tcpActionBuilder{}
}