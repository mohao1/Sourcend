package test

import (
	"Sourcend/command"
	"Sourcend/mutation"
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

	err = sourcendServer.RegisterCommandHandler("command-1", &commandHandlerTest{})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	err = sourcendServer.RegisterMutationHandler("afterMutation", "mutation-1", server.After, &afterHandlerTest{})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	err = sourcendServer.RegisterMutationHandler("beforeMutation", "mutation-1", server.Before, &beforeHandlerTest{})
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	info := server.SourcendInfo{
		CommandID:  "command-1",
		MutationID: "mutation-2",
		Event:      "测试数据",
		Params:     map[string]string{},
	}
	err = sourcendServer.Execute(ctx, info)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
}

type commandHandlerTest struct {
}

func (c *commandHandlerTest) Handler(ctx context.Context, data command.CommandData) (error, command.CommandData) {
	fmt.Println(data.Event)
	return nil, data
}

type afterHandlerTest struct {
}

func (a *afterHandlerTest) Handler(ctx context.Context, data mutation.MutationData) (error, mutation.MutationData) {
	fmt.Println("after:", data.Event)
	return nil, data
}

type beforeHandlerTest struct {
}

func (a *beforeHandlerTest) Handler(ctx context.Context, data mutation.MutationData) (error, mutation.MutationData) {
	fmt.Println("before:", data.Event)
	return nil, data
}
