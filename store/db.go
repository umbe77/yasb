package store

import (
	"context"

	"github.com/umbe77/yasb/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	db_name string = "yasb"
)

type Storage struct {
	Client *mongo.Client
}

func NewStore(uri string) (*Storage, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return &Storage{
		Client: client,
	}, nil
}

func (s *Storage) getCollection(coll string) *mongo.Collection {
	return s.Client.Database(db_name).Collection(coll)
}

func (s *Storage) GetWorflows(ctx context.Context, json_filter string) ([]models.Workflow, error) {
	var filter bson.D
	if err := bson.UnmarshalExtJSON([]byte(json_filter), true, &filter); err != nil {
		return nil, err
	}

	coll := s.getCollection("workflows")
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	results := make([]models.Workflow, 0)
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *Storage) AddWorkflow(ctx context.Context, wf models.Workflow) error {
	coll := s.getCollection("workflow")
	if _, err := coll.InsertOne(ctx, wf); err != nil {
		return err
	}
	return nil
}
