package builders

import (
	"github.com/hellgate75/go-network/model"
	"github.com/hellgate75/go-network/model/context"
)

type ApiActionBuilder interface {
	With(function model.ApiActionFunction) ApiActionBuilder
	Build() model.ApiAction
}

type apiActionBuilder struct {
	function			model.ApiActionFunction
}

func (builder *apiActionBuilder) With(function model.ApiActionFunction) ApiActionBuilder {
	builder.function = function
	return builder
}

func (builder *apiActionBuilder) Build() model.ApiAction {
	return &apiAction{
		function: builder.function,
	}
}

type apiAction struct{
	function model.ApiActionFunction
	context *context.ApiCallContext
}

func (action *apiAction) With(context context.ApiCallContext) model.ApiAction {
	action.context = &context
	return action
}

func (action *apiAction) Do() error {
	return action.function(*action.context)
}

func NewApiActionBuilder() ApiActionBuilder {
	return &apiActionBuilder{}
}