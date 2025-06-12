package server

// 解析配置结构

// CommandConfig Command配置解析结构
type CommandConfig struct {
	CommandName           string                 `yaml:"CommandName"`           // Command名称
	AfterMutation         []string               `yaml:"AfterMutation"`         // Mutation顺序
	BeforeMutation        []string               `yaml:"BeforeMutation"`        // Mutation顺序
	StoreEvents           []string               `yaml:"StoreEvents"`           // StoreEvents列表
	CommandHandlerConfigs []CommandHandlerConfig `yaml:"CommandHandlerConfigs"` // CommandHandler配置文件
}

// CommandHandlerConfig CommandHandler配置文件
type CommandHandlerConfig struct {
	CommandId string            `yaml:"CommandID"`
	Params    map[string]string `yaml:"Params"` // Params - 扩展数据
}

// MutationType Mutation类型
type MutationType string

const (
	After  = "after"
	Before = "before"
)

// MutationConfig Mutation配置解析结构
type MutationConfig struct {
	MutationName           string                  `yaml:"MutationName"`           // Mutation名称
	MutationType           MutationType            `yaml:"MutationType"`           // Mutation类型
	MutationHandlerConfigs []MutationHandlerConfig `yaml:"MutationHandlerConfigs"` //MutationHandler配置文件
}

// MutationHandlerConfig MutationHandler配置文件
type MutationHandlerConfig struct {
	MutationId string            `yaml:"MutationId"`
	Params     map[string]string `yaml:"Params"` // Params - 扩展数据
}

// 解析Mutation配置文件夹的函数
func mutationConfigYamlDir(dir string) ([]MutationConfig, []MutationConfig, error) {
	return nil, nil, nil
}

// 解析Command配置文件的函数
func commandConfigYamlDir(path string) (CommandConfig, error) {

	return CommandConfig{}, nil
}
