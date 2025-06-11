package server

import (
	"Sourcend/command"
	"Sourcend/common"
	"Sourcend/mutation"
	"Sourcend/store_event"
	"context"
	"errors"
)

type SourcendInfo struct {
	CommandID  string            // CommandID - CommandID指令
	MutationID string            // MutationID - MutationID指令
	Event      string            // Event - 修改数据 - 序列化成JSON字符串的结构
	Params     map[string]string // Params - 其他扩展数据
}

// SourcendManager Command-Mutation指令的管理器
type sourcendManager struct {
	commandManager  *command.Manager         // command的管理器
	afterMutations  []*mutation.Manager      // afterMutation的管理器
	beforeMutations []*mutation.Manager      // beforeMutation的管理器
	storeEvents     []store_event.StoreEvent // StoreEvent列表
}

// Execute 执行CommandExecute操作
func (s *sourcendManager) Execute(ctx context.Context, info SourcendInfo) error {
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

// RegisterCommand 注册CommandManager
func (s *sourcendManager) RegisterCommand(commandManager *command.Manager) error {
	if commandManager == nil {
		return errors.New("commandManager is nil")
	}
	s.commandManager = commandManager
	return nil
}

// RegisterStoreEvent 注册StoreEvent
func (s *sourcendManager) RegisterStoreEvent(storeEvent store_event.StoreEvent) error {
	s.storeEvents = append(s.storeEvents, storeEvent)
	return nil
}

// RegisterAfterMutations 注册AfterMutationManager
func (s *sourcendManager) RegisterAfterMutations(mutationManager *mutation.Manager) error {
	if mutationManager == nil {
		return errors.New("afterMutationManager is nil")
	}
	s.afterMutations = append(s.afterMutations, mutationManager)
	return nil
}

// RegisterBeforeMutations 注册BeforeMutationManager
func (s *sourcendManager) RegisterBeforeMutations(mutationManager *mutation.Manager) error {
	if mutationManager == nil {
		return errors.New("beforeMutationManager is nil")
	}
	s.beforeMutations = append(s.beforeMutations, mutationManager)
	return nil
}
