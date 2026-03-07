package main

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
)

func NewChatModel(ctx context.Context) *ark.ChatModel {
	model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: os.Getenv("ARK_API_KEY"),
		Model:  os.Getenv("MODEL"),
	})
	if err != nil {
		panic(err)
	}
	return model
}
