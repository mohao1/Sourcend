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

type MySQLConfig struct {
	ActionTable  string // Action存储表名
	IsActionKey  bool   // 是否开启Key的配置模式
	ActionKey    string // Action的Key
	ActionMaxLen int64  // Action内的最大长度
	DSN          string // 连接的数据库配置
}
