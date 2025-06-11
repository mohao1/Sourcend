package server

import (
	"Sourcend/command"
	"Sourcend/mutation"
	"Sourcend/store_event"
)

// SourcendServer 服务主体
type SourcendServer struct {
	sourcendManager ManagerInterface
	// 配置文件管理模块
}

func NewDefaultSourcend() *SourcendServer {

	// 生成Mutation的操作
	// 生成Command的操作

	return &SourcendServer{sourcendManager: &sourcendManager{
		commandManager:  command.NewManager(),
		afterMutations:  make([]*mutation.Manager, 0),
		beforeMutations: make([]*mutation.Manager, 0),
		storeEvents:     make([]store_event.StoreEvent, 0),
	}}
}

func (s *SourcendServer) CommandUse() {

}
