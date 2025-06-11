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
	Execute(context.Context, SourcendInfo) error
	RegisterCommand(*command.Manager) error
	RegisterStoreEvent(store_event.StoreEvent) error
	RegisterAfterMutations(*mutation.Manager) error
	RegisterBeforeMutations(*mutation.Manager) error
}
