package command

import (
	"Sourcend/common"
	"context"
)

// Middleware 中间件的配置
type Middleware func(Handler) Handler

type Handler func(context.Context, CommandData) (error, CommandData)

// HandlerInterface MutationHandler接口
type HandlerInterface interface {
	Handler(ctx context.Context, data CommandData) (error, CommandData)
}

type StoreEvent interface {
	Handler(ctx context.Context, data CommandData) error
}

// Interface Manager接口
type Interface interface {
	// Use 注册Command的中间件,全局的执行
	Use(...Middleware) error
	// apply 配置Command指令
	apply(string, Handler) (Handler, error)
	//Register 注册Command指令
	Register(string, HandlerInterface, ...Middleware) error

	// 执行Command的Execute操作
	executeCommand(context.Context, CommandData) (CommandData, error)

	// 设置存储Event的记录信息
	storeEvent(context.Context, CommandData) error

	// UseStoreEvent 配置存储Event的自定义的操作
	UseStoreEvent(StoreEvent) error

	// Execute 对外执行Execute
	Execute(context.Context, common.CommandInfo) (*common.CommandInfo, error)
}
