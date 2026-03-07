package main

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus"
)

func NewArkRetriever(ctx context.Context, embedder *ark.Embedder) *milvus.Retriever {
	collectName := "test"
	retriever, err := milvus.NewRetriever(ctx, &milvus.RetrieverConfig{
		Client:       MilvusCli,
		Collection:   collectName,
		VectorField:  "vector",
		OutputFields: []string{"id", "content", "metadata"},
		TopK:         1,
		Embedding:    embedder,
	})
	if err != nil {
		panic(err)
	}
	return retriever
}
