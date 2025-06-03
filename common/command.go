package common

// CommandInfo Command层的数据接收结构
type CommandInfo struct {
	CommandID  string            // CommandID - CommandID指令
	MutationID string            // MutationID - MutationID指令
	Event      string            // Event - Command携带的数据 - 序列化成JSON字符串的结构
	Params     map[string]string // Params - 其他扩展数据
}
