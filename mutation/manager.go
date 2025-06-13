package mutation

import (
	"Sourcend/common"
	"context"
	"errors"
	"fmt"
	"sync"
)

// Manager Mutation的管理器
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

// Use 注册中间件的操作
func (m *Manager) Use(middlewares ...Middleware) error {
	m.middlewaresLock.Lock()
	defer m.middlewaresLock.Unlock()
	m.middlewares = append(m.middlewares, middlewares...)
	return nil
}

// apply加载中间件的操作
func (m *Manager) apply(mutationID string, handler Handler) (Handler, error) {
	m.middlewaresLock.Lock()
	defer m.middlewaresLock.Unlock()
	m.handlerMiddlewaresLock.Lock()
	defer m.handlerMiddlewaresLock.Unlock()

	if middlewares, ok := m.handlerMiddlewares[mutationID]; ok {
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

// Register 注册mutation
func (m *Manager) Register(mutationID string, handler HandlerInterface, middlewares ...Middleware) error {
	if mutationID == "" {
		return errors.New("mutationID must not be empty")
	}
	if len(middlewares) > 0 {
		m.handlerMiddlewaresLock.Lock()
		m.handlerMiddlewares[mutationID] = middlewares
		m.handlerMiddlewaresLock.Unlock()
	}
	// 加载handler的中间件
	apply, err := m.apply(mutationID, handler.Handler)
	if err != nil {
		return err
	}

	// 注册进入handler
	m.handlerManager[mutationID] = apply
	return nil
}

// Execute 执行
func (m *Manager) Execute(ctx context.Context, data common.MutationInfo) error {
	mutationID := data.MutationID
	mutationHandler, ok := m.handlerManager[mutationID]

	if !ok {
		return errors.New(fmt.Sprintf("mutation handler not found: %s", mutationID))
	}

	// 转换数据 Info=>Data
	config := m.ManagerConfig.MutationConfigMap[data.MutationID]
	mutationData := MutationData{
		MutationID: data.MutationID,
		Event:      data.Event,
		Params:     data.Params,
		Config:     config,
	}

	// 执行函数
	err, mutationData := mutationHandler(ctx, mutationData)
	if err != nil {
		return err
	}

	return nil
}
