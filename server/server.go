package server

import (
	"Sourcend/command"
	"Sourcend/mutation"
	"Sourcend/store_event"
	"context"
	"errors"
	"fmt"
	"sync"
)

// SourcendServer 服务主体
type SourcendServer struct {
	sourcendManager ManagerInterface
	// 配置文件管理模块
	commandConfig   CommandConfig
	afterMutations  []MutationConfig
	beforeMutations []MutationConfig
}

// NewDefaultSourcend 根据配置文件生成SourcendServer
func NewDefaultSourcend(mutationYamlDir, commandYamlPath string) (*SourcendServer, error) {
	// 解析Yaml
	after, before, err := mutationConfigYamlDir(mutationYamlDir)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	commandConfig, err := commandConfigYamlPath(commandYamlPath)
	if err != nil {
		return nil, err
	}

	// 存储数据结构
	afterMutations := make([]*mutation.Manager, len(after))
	beforeMutations := make([]*mutation.Manager, len(before))
	storeEvents := make([]store_event.StoreEvent, len(commandConfig.StoreEvents))
	var commandManager *command.Manager

	group := sync.WaitGroup{}
	group.Add(4)

	// 生成Command的操作
	// 处理command的协程
	go func() {
		defer group.Done()
		commandConfigMap := make(map[string]command.CommandConfig)
		for _, c := range commandConfig.CommandHandlerConfigs {
			commandConfigMap[c.CommandId] = command.CommandConfig{
				CommandID: c.CommandId,
				Params:    c.Params,
			}
		}
		managerConfig := command.ManagerConfig{
			ManagerName:      commandConfig.CommandName,
			CommandConfigMap: commandConfigMap,
		}
		commandManager = command.NewManager(managerConfig)
	}()

	// 生成Mutation的操作
	// 处理after的协程
	go func() {
		defer group.Done()
		// 排序
		for k, config := range after {
			mutationConfigMap := make(map[string]mutation.MutationConfig)
			for _, c := range config.MutationHandlerConfigs {
				mutationConfigMap[c.MutationId] = mutation.MutationConfig{
					MutationId: c.MutationId,
					Params:     c.Params,
				}
			}
			managerConfig := mutation.ManagerConfig{
				ManagerName:       config.MutationName,
				MutationConfigMap: mutationConfigMap,
			}
			afterMutations[k] = mutation.NewManager(managerConfig)
		}
	}()

	// 处理before的协程
	go func() {
		defer group.Done()
		// 排序
		for k, config := range before {
			mutationConfigMap := make(map[string]mutation.MutationConfig)
			for _, c := range config.MutationHandlerConfigs {
				mutationConfigMap[c.MutationId] = mutation.MutationConfig{
					MutationId: c.MutationId,
					Params:     c.Params,
				}
			}
			managerConfig := mutation.ManagerConfig{
				ManagerName:       config.MutationName,
				MutationConfigMap: mutationConfigMap,
			}
			beforeMutations[k] = mutation.NewManager(managerConfig)
		}
	}()

	// 生成StoreEvent的操作
	go func() {
		defer group.Done()
		if commandConfig.StoreEvents == nil {
			return
		}
		storeEventList := commandConfig.StoreEvents
		for _, store := range storeEventList {
			switch store {
			case store_event.MySQL:
			case store_event.Redis:
			default:
				fmt.Println("store does not match")
			}
		}
	}()

	group.Wait()
	return &SourcendServer{
		sourcendManager: &sourcendManager{
			commandManager:  commandManager,
			afterMutations:  afterMutations,
			beforeMutations: beforeMutations,
			storeEvents:     storeEvents,
		},
		commandConfig:   *commandConfig,
		afterMutations:  after,
		beforeMutations: before,
	}, nil
}

// CommandUse Command注册对应的拦截器
func (s *SourcendServer) CommandUse(middlewares ...command.Middleware) error {
	err := s.sourcendManager.CommandUse(middlewares...)
	if err != nil {
		return err
	}
	return nil
}

// RegisterCommandHandler 注册CommandHandler
func (s *SourcendServer) RegisterCommandHandler(commandID string, handler command.HandlerInterface, middlewares ...command.Middleware) error {
	err := s.sourcendManager.RegisterCommandHandler(commandID, handler, middlewares...)
	if err != nil {
		return err
	}
	return nil
}

// MutationUse 对应的Mutation注册对应的拦截器
func (s *SourcendServer) MutationUse(mutationName string, mutationType MutationType, middlewares ...mutation.Middleware) error {
	switch mutationType {
	case After:
		err := s.sourcendManager.AfterMutationUse(mutationName, middlewares...)
		if err != nil {
			return err
		}
	case Before:
		err := s.sourcendManager.BeforeMutationUse(mutationName, middlewares...)
		if err != nil {
			return err
		}
	default:
		return errors.New("mutation type not supported")
	}
	return nil
}

// RegisterMutationHandler 注册MutationHandler
func (s *SourcendServer) RegisterMutationHandler(mutationName string, mutationID string, mutationType MutationType, handler mutation.HandlerInterface, middlewares ...mutation.Middleware) error {
	switch mutationType {
	case After:
		err := s.sourcendManager.RegisterAfterMutationHandler(mutationName, mutationID, handler, middlewares...)
		if err != nil {
			return err
		}
	case Before:
		err := s.sourcendManager.RegisterBeforeMutationHandler(mutationName, mutationID, handler, middlewares...)
		if err != nil {
			return err
		}
	default:
		return errors.New("mutation type not supported")
	}
	return nil
}

// RegisterStoreEvents 手动注入自定义的StoreEvent
func (s *SourcendServer) RegisterStoreEvents(storeEvent store_event.StoreEvent) error {
	if storeEvent == nil {
		err := errors.New("store event is nil")
		fmt.Println(err)
		return err
	}
	err := s.sourcendManager.RegisterStoreEvent(storeEvent)
	if err != nil {
		return err
	}
	return nil
}

// Execute 执行Execute
func (s *SourcendServer) Execute(ctx context.Context, info SourcendInfo) error {
	err := s.sourcendManager.Execute(ctx, info)
	if err != nil {
		return err
	}
	return nil
}
