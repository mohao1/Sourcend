package mutation

import (
	"Sourcend/common"
	"context"
)

// Middleware 中间件的配置
type Middleware func(Handler) Handler

type Handler func(context.Context, MutationData) (error, MutationData)

// HandlerInterface MutationHandler接口
type HandlerInterface interface {
	Handler(ctx context.Context, data MutationData) (error, MutationData)
}

// Interface Manager接口
type Interface interface {
	// Use 注册Mutation的中间件,全局的执行
	Use(...Middleware) error
	// apply 配置Mutation指令
	apply(string, Handler) (Handler, error)
	//Register 注册Mutation指令
	Register(string, HandlerInterface, ...Middleware) error
	// Execute 执行Execute
	Execute(context.Context, common.MutationInfo) (*common.MutationInfo, error)
}
