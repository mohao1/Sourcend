package command

import (
	"Sourcend/common"
	"Sourcend/mutation"
	"Sourcend/store_event"
	"context"
)

// Middleware 中间件的配置
type Middleware func(Handler) Handler

type Handler func(context.Context, CommandData) (error, CommandData)

// HandlerInterface MutationHandler接口
type HandlerInterface interface {
	Handler(ctx context.Context, data CommandData) (error, CommandData)
}

// Interface Manager接口
// executeCommand => beforeMutation => storeEvent => afterMutation
type Interface interface {
	// Use 注册Command的中间件,全局的执行
	Use(...Middleware) error
	// apply 配置Command指令
	apply(string, Handler) (Handler, error)
	//Register 注册Command指令
	Register(string, HandlerInterface, ...Middleware) error

	// 执行Command的Execute操作
	executeCommand(context.Context, common.CommandInfo) (*common.CommandInfo, error)
	// Execute 对外执行Execute
	Execute(context.Context, common.CommandInfo, []store_event.StoreEvent, []*mutation.Manager, []*mutation.Manager) error

	// StoreEvent之前执行
	beforeMutation(context.Context, common.CommandInfo, []*mutation.Manager) error
	// StoreEvent之后执行
	afterMutation(context.Context, common.CommandInfo, []*mutation.Manager) error

	// 存储Event的记录信息操作
	storeEvent(context.Context, common.CommandInfo, []store_event.StoreEvent)
}
