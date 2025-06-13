package server

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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

// 检查文件是否为YAML文件
func isYAMLFile(fileName string) bool {
	lowerName := strings.ToLower(fileName)
	return strings.HasSuffix(lowerName, ".yaml") || strings.HasSuffix(lowerName, ".yml")
}

// 遍历目录获取所有YAML文件
func getYAMLFilesInDir(dirPath string) ([]string, error) {
	var yamlFiles []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isYAMLFile(info.Name()) {
			yamlFiles = append(yamlFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历目录 %s 失败: %v", dirPath, err)
	}

	return yamlFiles, nil
}

// 读取YAMLMutationConfig文件内容
func readYAMLFileByMutationConfig(filePath string) (*MutationConfig, error) {
	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件 %s 失败: %v", filePath, err)
	}

	// 解析YAML内容
	var config MutationConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析YAML文件 %s 失败: %v", filePath, err)
	}

	return &config, nil
}

// 读取YAMLCommandConfig文件内容
func readYAMLFileByCommandConfig(filePath string) (*CommandConfig, error) {
	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取文件 %s 失败: %v", filePath, err)
	}

	// 解析YAML内容
	var config CommandConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析YAML文件 %s 失败: %v", filePath, err)
	}

	return &config, nil
}

// 解析Mutation配置文件夹的函数
func mutationConfigYamlDir(dir string) ([]MutationConfig, []MutationConfig, error) {
	afterList := make([]MutationConfig, 0)
	beforeList := make([]MutationConfig, 0)

	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("错误: 目录 %s 不存在\n", dir)
		return nil, nil, errors.New(fmt.Sprintf("错误: 目录 %s 不存在\n", dir))
	}

	// 获取所有YAML文件
	yamlFiles, err := getYAMLFilesInDir(dir)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return nil, nil, errors.New(fmt.Sprintf("错误: %v\n", err))
	}

	if len(yamlFiles) == 0 {
		fmt.Printf("目录 %s 中没有找到YAML文件\n", dir)
		return nil, nil, errors.New(fmt.Sprintf("目录 %s 中没有找到YAML文件\n", dir))
	}

	// 读取并解析每个YAML文件
	for _, file := range yamlFiles {
		config, err := readYAMLFileByMutationConfig(file)
		if err != nil {
			fmt.Printf("警告: %v\n", err)
			continue
		}

		// 初始化成长度为0列表
		if config.MutationHandlerConfigs == nil {
			config.MutationHandlerConfigs = make([]MutationHandlerConfig, 0)
		}

		// 加载数据进入列表
		switch config.MutationType {
		case After:
			afterList = append(afterList, *config)
		case Before:
			beforeList = append(beforeList, *config)
		default:
			return nil, nil, errors.New("MutationType is error")
		}
	}

	return afterList, beforeList, nil
}

// 解析Command配置文件的函数
func commandConfigYamlPath(path string) (*CommandConfig, error) {
	config, err := readYAMLFileByCommandConfig(path)
	if err != nil {
		fmt.Printf("警告: %v\n", err)
		return nil, err
	}
	if config.CommandHandlerConfigs == nil {
		config.CommandHandlerConfigs = make([]CommandHandlerConfig, 0)
	}

	if config.StoreEvents == nil {
		config.StoreEvents = make([]string, 0)
	}

	if config.AfterMutation == nil {
		config.AfterMutation = make([]string, 0)
	}

	if config.BeforeMutation == nil {
		config.BeforeMutation = make([]string, 0)
	}

	return config, nil
}
