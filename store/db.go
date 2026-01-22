package store

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Storage struct {
	ctx    context.Context
	Client *mongo.Client
}

func NewStore(ctx context.Context, uri string) (*Storage, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &Storage{
		ctx:    ctx,
		Client: client,
	}, nil
}
