package common

// MutationInfo Command => Mutation 传输数据结构
type MutationInfo struct {
	MutationID string            // MutationID - MutationID指令
	Event      string            // Event - Mutation携带的数据 - 序列化成JSON字符串的结构
	Params     map[string]string // Params - 其他扩展数据
}
