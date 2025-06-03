package server

import (
	"Sourcend/command"
	"Sourcend/common"
	"Sourcend/mutation"
	"Sourcend/store_event"
	"context"
)

type SourcendInfo struct {
	CommandID  string            // CommandID - CommandID指令
	MutationID string            // MutationID - MutationID指令
	Event      string            // Event - 修改数据 - 序列化成JSON字符串的结构
	Params     map[string]string // Params - 其他扩展数据
}

// SourcendManager Command-Mutation指令的管理器
type SourcendManager struct {
	commandManager  command.Manager          // command的管理器
	afterMutations  []*mutation.Manager      // afterMutation的管理器
	beforeMutations []*mutation.Manager      // beforeMutation的管理器
	storeEvents     []store_event.StoreEvent // StoreEvent列表
}

func (s *SourcendManager) Execute(ctx context.Context, info SourcendInfo) error {
	// 执行commandManager指令
	data := common.CommandInfo{
		CommandID:  info.CommandID,
		MutationID: info.MutationID,
		Event:      info.Event,
		Params:     info.Params,
	}
	err := s.commandManager.Execute(ctx, data, s.storeEvents, s.beforeMutations, s.afterMutations)
	if err != nil {
		return err
	}
	return nil
}
