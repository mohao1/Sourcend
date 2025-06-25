package store_event

type StoreEventInfo struct {
	CommandID  string            // CommandID - CommandID指令
	MutationID string            // MutationID - MutationID指令
	Event      string            // Event - 修改数据 - 序列化成JSON字符串的结构
	Params     map[string]string // Params - 其他扩展数据
}

type StoreEventType string

const (
	MySQL = "mysql"
	Redis = "redis"
)
