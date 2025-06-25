package command

import (
	"Sourcend/common"
	"Sourcend/mutation"
	"Sourcend/store_event"
	"context"
	"errors"
	"fmt"
	"sync"
)

// Manager Command的管理器
type Manager struct {
	ManagerName            string                  // 管理器的名称
	ManagerConfig          ManagerConfig           // 管理器的配置文件
	middlewares            []Middleware            // 中间件的存储
	middlewaresLock        sync.RWMutex            // 中间件的存储的锁
	handlerMiddlewares     map[string][]Middleware // Handler自己的Middlewares
	handlerMiddlewaresLock sync.RWMutex            // Handler自己的Middlewares间件的存储的锁
	handlerManager         map[string]Handler      // Handler存储
}

func NewManager(config ManagerConfig) *Manager {
	return &Manager{
		ManagerName:        config.ManagerName,
		ManagerConfig:      config,
		middlewares:        make([]Middleware, 0),
		handlerMiddlewares: make(map[string][]Middleware),
		handlerManager:     make(map[string]Handler),
	}
}

func (m *Manager) Use(middlewares ...Middleware) error {
	m.middlewaresLock.Lock()
	defer m.middlewaresLock.Unlock()
	m.middlewares = append(m.middlewares, middlewares...)
	return nil
}

func (m *Manager) apply(commandID string, handler Handler) (Handler, error) {
	m.middlewaresLock.Lock()
	defer m.middlewaresLock.Unlock()
	m.handlerMiddlewaresLock.Lock()
	defer m.handlerMiddlewaresLock.Unlock()

	if middlewares, ok := m.handlerMiddlewares[commandID]; ok {
		for i := len(middlewares) - 1; i >= 0; i-- {
			middleware := middlewares[i]
			handler = middleware(handler)
		}
	}

	for i := len(m.middlewares) - 1; i >= 0; i-- {
		middleware := m.middlewares[i]
		handler = middleware(handler)
	}

	return handler, nil
}

func (m *Manager) Register(commandID string, handler HandlerInterface, middlewares ...Middleware) error {
	if commandID == "" {
		return errors.New("mutationID must not be empty")
	}
	if len(middlewares) > 0 {
		m.handlerMiddlewaresLock.Lock()
		m.handlerMiddlewares[commandID] = middlewares
		m.handlerMiddlewaresLock.Unlock()
	}
	// 加载handler的中间件
	apply, err := m.apply(commandID, handler.Handler)
	if err != nil {
		return err
	}

	// 注册进入handler
	m.handlerManager[commandID] = apply
	return nil
}

func (m *Manager) executeCommand(ctx context.Context, data common.CommandInfo) (*common.CommandInfo, error) {
	commandID := data.CommandID
	commandHandler, ok := m.handlerManager[commandID]

	if !ok {
		return nil, errors.New("command handler not found: " + commandID)
	}

	// 转换数据 Info=>Data
	config := m.ManagerConfig.CommandConfigMap[commandID]
	commandData := CommandData{
		CommandID:  commandID,
		MutationID: data.MutationID,
		Event:      data.Event,
		Params:     data.Params,
		Config:     config,
	}

	// 执行函数
	err, commandData := commandHandler(ctx, commandData)
	if err != nil {
		return nil, err
	}

	// 转换数据 CommandData=>MutationInfo
	newData := &common.CommandInfo{
		CommandID:  commandData.CommandID,
		MutationID: commandData.MutationID,
		Event:      commandData.Event,
		Params:     commandData.Params,
	}
	return newData, nil
}

func (m *Manager) Execute(ctx context.Context, info common.CommandInfo, storeEvents []store_event.StoreEvent, beforeMutations []*mutation.Manager, afterMutations []*mutation.Manager) error {

	// 执行Command
	commandInfo, err := m.executeCommand(ctx, info)
	if err != nil {
		return err
	}

	// 执行beforeMutation
	err = m.beforeMutation(ctx, *commandInfo, beforeMutations)
	if err != nil {
		return err
	}

	//执行StoreEvent
	m.storeEvent(ctx, *commandInfo, storeEvents)

	// 执行afterMutation
	err = m.afterMutation(ctx, *commandInfo, afterMutations)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) beforeMutation(ctx context.Context, info common.CommandInfo, managers []*mutation.Manager) error {
	data := common.MutationInfo{
		MutationID: info.MutationID,
		Event:      info.Event,
		Params:     info.Params,
	}
	for _, manager := range managers {
		err := manager.Execute(ctx, data)
		if err != nil {
			fmt.Println("mutation err:", err)
			continue
		}
	}
	return nil
}

func (m *Manager) afterMutation(ctx context.Context, info common.CommandInfo, managers []*mutation.Manager) error {
	data := common.MutationInfo{
		MutationID: info.MutationID,
		Event:      info.Event,
		Params:     info.Params,
	}
	for _, manager := range managers {
		err := manager.Execute(ctx, data)
		if err != nil {
			fmt.Println("mutation err:", err)
			continue
		}
	}
	return nil
}

func (m *Manager) storeEvent(ctx context.Context, info common.CommandInfo, storeEvents []store_event.StoreEvent) {
	data := store_event.StoreEventInfo{
		CommandID:  info.CommandID,
		MutationID: info.MutationID,
		Event:      info.Event,
		Params:     info.Params,
	}
	for _, storeEvent := range storeEvents {
		err := storeEvent.Handler(ctx, data)
		if err != nil {
			fmt.Println("store err:", err)
			continue
		}
	}
}
