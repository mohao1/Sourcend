package command

import "sync"

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
