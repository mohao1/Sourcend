package server

import (
	"Sourcend/command"
	"Sourcend/common"
	"Sourcend/mutation"
	"Sourcend/store_event"
	"context"
	"errors"
	"fmt"
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

// CommandUse 注册Command的中间件
func (s *sourcendManager) CommandUse(middlewares ...command.Middleware) error {
	err := s.commandManager.Use(middlewares...)
	if err != nil {
		fmt.Println("commandManager Use err:", err)
		return err
	}
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

// AfterMutationUse 注册AfterMutation的中间件
func (s *sourcendManager) AfterMutationUse(mutationName string, middlewares ...mutation.Middleware) error {
	for _, m := range s.afterMutations {
		if m.ManagerName == mutationName {
			err := m.Use(middlewares...)
			if err != nil {
				fmt.Println("AfterMutation Use err:", err)
				return err
			}
			return nil
		}
	}
	return errors.New(fmt.Sprintf("mutation %s not found in after mutations", mutationName))
}

// BeforeMutationUse 注册BeforeMutation的中间件
func (s *sourcendManager) BeforeMutationUse(mutationName string, middlewares ...mutation.Middleware) error {
	for _, m := range s.beforeMutations {
		if m.ManagerName == mutationName {
			err := m.Use(middlewares...)
			if err != nil {
				fmt.Println("BeforeMutation Use err:", err)
				return err
			}
			return nil
		}
	}
	return errors.New(fmt.Sprintf("mutation %s not found in after mutations", mutationName))
}

// RegisterAfterMutationHandler 往AfterMutation中注入Handler的方法
func (s *sourcendManager) RegisterAfterMutationHandler(mutationName string, mutationID string, handler mutation.HandlerInterface, middlewares ...mutation.Middleware) error {
	for _, m := range s.afterMutations {
		if m.ManagerName == mutationName {
			err := m.Register(mutationID, handler, middlewares...)
			if err != nil {
				fmt.Println("AfterMutation Register err:", err)
				return err
			}
		}
	}
	return nil
}

// RegisterBeforeMutationHandler 往BeforeMutation中注入Handler的方法
func (s *sourcendManager) RegisterBeforeMutationHandler(mutationName string, mutationID string, handler mutation.HandlerInterface, middlewares ...mutation.Middleware) error {
	for _, m := range s.beforeMutations {
		if m.ManagerName == mutationName {
			err := m.Register(mutationID, handler, middlewares...)
			if err != nil {
				fmt.Println("BeforeMutation Register err:", err)
				return err
			}
		}
	}
	return nil
}

// RegisterCommandHandler 往Command中注入Handler方法
func (s *sourcendManager) RegisterCommandHandler(commandID string, handler command.HandlerInterface, middlewares ...command.Middleware) error {
	err := s.commandManager.Register(commandID, handler, middlewares...)
	if err != nil {
		fmt.Println("Command Register err:", err)
		return err
	}
	return nil
}
