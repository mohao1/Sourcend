package test

import (
	"Sourcend/server"
	"context"
	"fmt"
	"testing"
)

func TestSourcendServer(t *testing.T) {
	ctx := context.Background()
	sourcendServer, err := server.NewDefaultSourcend(
		"./config/mutation_config",
		"./config/command_config/command.yaml",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	info := server.SourcendInfo{
		CommandID:  "command-1",
		MutationID: "mutation-1",
		Event:      "测试数据",
		Params:     map[string]string{},
	}
	err = sourcendServer.Execute(ctx, info)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
}
