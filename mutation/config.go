package mutation

// MutationConfig 各个Mutation配置
type MutationConfig struct {
	CommandID string
	Params    map[string]string // Params - 扩展数据
}

// MutationData Mutation层接收的数据
type MutationData struct {
	MutationID string            // MutationID - MutationID指令
	Event      string            // Event - Mutation携带的数据 - 序列化成JSON字符串的结构
	Params     map[string]string // Params - 其他扩展数据
	Config     MutationConfig    // 配置文件
}

// ManagerConfig MutationManager的配置文件的存储
type ManagerConfig struct {
	ManagerName       string
	MutationConfigMap map[string]MutationConfig
}
