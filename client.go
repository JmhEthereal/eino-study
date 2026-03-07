package main

import (
	"context"

	cli "github.com/milvus-io/milvus-sdk-go/v2/client"
)

var MilvusCli cli.Client

func InitClint() {
	ctx := context.Background()
	client, err := cli.NewClient(ctx, cli.Config{
		Address: "localhost:19530",
		DBName:  "my_eino",
	})
	if err != nil {
		panic(err)
	}
	MilvusCli = client
}
