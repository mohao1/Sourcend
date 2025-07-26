package action

// Action action的根数据 - 保存数据结构
type Action struct {
	ActionName   string  // Action的名称
	RootData     string  // 初始化的数据 - json数据
	StartEventID int64   // 起始版本ID
	Events       []Event // 指令版本记录 DATA:Event数据
	EventLen     int64   // 记录长度
}
