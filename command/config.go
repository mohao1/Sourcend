package command

// CommandConfig 各个Command配置
type CommandConfig struct {
	CommandID string
	Params    map[string]string // Params - 扩展数据
}

// CommandData Command层接收的数据
type CommandData struct {
	CommandID  string            // CommandID - CommandID指令
	MutationID string            // MutationID - MutationID指令
	Event      string            // Event - Mutation携带的数据 - 序列化成JSON字符串的结构
	Params     map[string]string // Params - 其他扩展数据
	Config     CommandConfig     // 配置文件
}

// ManagerConfig CommandManager的配置文件的存储
type ManagerConfig struct {
	ManagerName      string
	CommandConfigMap map[string]CommandConfig
}
