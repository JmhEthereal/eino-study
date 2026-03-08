package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"os"
)

type State struct {
	History map[string]any
}

func genFunc(ctx context.Context) *State {
	return &State{
		History: make(map[string]any),
	}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	g := compose.NewGraph[map[string]string, *schema.Message](
		compose.WithGenLocalState(genFunc))
	lambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output map[string]string, err error) {
		_ = compose.ProcessState[*State](ctx, func(_ context.Context, state *State) error {
			state.History["tsundere_action"] = "我喜欢你"
			state.History["cute_action"] = "摸摸头"
			return nil
		})
		if input["role"] == "tsundere" {
			return map[string]string{"role": "傲娇的", "content": input["content"]}, nil
		} else if input["role"] == "cute" {
			return map[string]string{"role": "可爱的", "content": input["content"]}, nil
		}
		return map[string]string{"role": input["role"], "content": input["content"]}, nil
	})
	TsundereLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		_ = compose.ProcessState[*State](ctx, func(_ context.Context, state *State) error {
			input["content"] = input["content"] + state.History["tsundere_action"].(string)
			return nil
		})
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个高冷傲娇的大小姐，每次都会用傲娇的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	CuteLambda := compose.InvokableLambda(func(ctx context.Context, input map[string]string) (output []*schema.Message, err error) {
		return []*schema.Message{
			{
				Role:    schema.System,
				Content: "你是一个可爱的小女孩，每次都会用可爱的语气回答我的问题",
			},
			{
				Role:    schema.User,
				Content: input["content"],
			},
		}, nil
	})
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("MODEL"),
	})
	if err != nil {
		panic(err)
	}
	cutePreHandler := func(ctx context.Context, input map[string]string, state *State) (map[string]string, error) {
		input["content"] = input["content"] + state.History["cute_action"].(string)
		return input, nil
	}
	//注册节点
	err = g.AddLambdaNode("lambda", lambda)
	if err != nil {
		panic(err)
	}
	err = g.AddLambdaNode("tsundere", TsundereLambda)
	if err != nil {
		panic(err)
	}
	err = g.AddLambdaNode("cute", CuteLambda, compose.WithStatePreHandler(cutePreHandler))
	if err != nil {
		panic(err)
	}
	err = g.AddChatModelNode("model", model)
	if err != nil {
		panic(err)
	}
	//链接节点
	err = g.AddEdge(compose.START, "lambda")
	if err != nil {
		panic(err)
	}
	//加入分支
	g.AddBranch("lambda", compose.NewGraphBranch(func(ctx context.Context, in map[string]string) (endNode string, err error) {
		if in["role"] == "傲娇的" {
			return "tsundere", nil
		}
		if in["role"] == "可爱的" {
			return "cute", nil
		}
		return "傲娇的", nil
	}, map[string]bool{"tsundere": true, "cute": true}))
	err = g.AddEdge("tsundere", "model")
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("cute", "model")
	if err != nil {
		panic(err)
	}
	err = g.AddEdge("model", compose.END)
	if err != nil {
		panic(err)
	}
	//编译
	r, err := g.Compile(ctx)
	if err != nil {
		panic(err)
	}
	//执行
	answer, err := r.Invoke(ctx, map[string]string{"role": "cute", "content": "你好啊"}, compose.WithCallbacks(getCallBack()))
	if err != nil {
		panic(err)
	}
	fmt.Println(answer.Content)
}

func getCallBack() callbacks.Handler {
	handler := callbacks.NewHandlerBuilder().OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
		fmt.Printf("当前%s节点输入:%s\n", info.Component, input)
		return ctx
	}).OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
		fmt.Printf("当前%s节点输出:%s\n", info.Component, output)
		return ctx
	}).Build()
	return handler
}
