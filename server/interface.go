package server

import (
	"Sourcend/command"
	"Sourcend/mutation"
	"Sourcend/store_event"
	"context"
)

// Interface Sourcend的服务
type Interface interface {
}

// ManagerInterface Sourcend的管理器
type ManagerInterface interface {
	// Execute执行处理
	Execute(context.Context, SourcendInfo) error
	// Command处理
	RegisterCommand(*command.Manager) error
	CommandUse(...command.Middleware) error
	RegisterCommandHandler(string, command.HandlerInterface, ...command.Middleware) error
	// StoreEvent处理
	RegisterStoreEvent(store_event.StoreEvent) error
	// After处理
	RegisterAfterMutations(*mutation.Manager) error
	AfterMutationUse(string, ...mutation.Middleware) error
	RegisterAfterMutationHandler(string, string, mutation.HandlerInterface, ...mutation.Middleware) error
	// Before处理
	RegisterBeforeMutations(*mutation.Manager) error
	BeforeMutationUse(string, ...mutation.Middleware) error
	RegisterBeforeMutationHandler(string, string, mutation.HandlerInterface, ...mutation.Middleware) error
}
