package data

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Models struct {
	Logs *LogStore
}

func New(client *mongo.Client) *Models {
	return &Models{
		Logs: NewLogStore(client, "appdb", "logs"),
	}
}

type LogStore struct {
	coll *mongo.Collection
}

func NewLogStore(client *mongo.Client, dbName, collName string) *LogStore {
	return &LogStore{
		coll: client.Database(dbName).Collection(collName),
	}
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (s *LogStore) Insert(ctx context.Context, entry *LogEntry) error {

	now := time.Now().UTC()
	entry.CreatedAt = now
	entry.UpdatedAt = now

	if _, err := s.coll.InsertOne(ctx, entry); err != nil {
		return fmt.Errorf("LogStore.Insert: %w", err)
	}

	return nil
}

func (s *LogStore) All(ctx context.Context) ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	findOpts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := s.coll.Find(ctx, bson.D{}, findOpts)
	if err != nil {
		return nil, fmt.Errorf("LogStore.All find: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(&item)
		if err != nil {
			return nil, fmt.Errorf("LogStore.All decode: %w", err)
		}

		logs = append(logs, &item)
	}

	return logs, nil
}

func (s *LogStore) GetOneLog(ctx context.Context, id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	collection := s.coll


	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id %q: %w", id, err)
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("LogStore.GetOneLog: %w", err)
	}

	return &entry, nil
}

func (s *LogStore) DropCollection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	collection := s.coll

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (s *LogStore) Update(ctx context.Context, logs *LogEntry) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	collection := s.coll

	docID, err := primitive.ObjectIDFromHex(logs.ID)
	if err != nil{
		return nil, err
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id":docID}, bson.D{
		{
			"$set", bson.D{
				{"name", logs.Name},
				{"data", logs.Data},
				{"updated_at", time.Now()},
			},
		},
	})


	if err != nil{
		return nil, err
	}

	return result, nil
}
