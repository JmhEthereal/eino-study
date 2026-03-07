package main

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func NewArkIndexer(ctx context.Context, embedder *ark.Embedder) *milvus.Indexer {
	collectName := "test"
	var fields = []*entity.Field{
		{
			Name:     "id",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "32",
			},
			PrimaryKey: true,
		},
		{
			Name:     "vector", // 确保字段名匹配
			DataType: entity.FieldTypeBinaryVector,
			TypeParams: map[string]string{
				"dim": "65536",
			},
		},
		{
			Name:     "content",
			DataType: entity.FieldTypeVarChar,
			TypeParams: map[string]string{
				"max_length": "8192",
			},
		},
		{
			Name:     "metadata",
			DataType: entity.FieldTypeJSON,
		},
	}
	indexer, err := milvus.NewIndexer(ctx, &milvus.IndexerConfig{
		Client:     MilvusCli,
		Collection: collectName,
		Fields:     fields,
		Embedding:  embedder,
	})
	if err != nil {
		log.Fatalf("Failed to create indexer: %v", err)
	}

	return indexer
}
