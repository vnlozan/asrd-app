package repo

import (
	"context"
	"log"
	"logger/internal/dto"
	"logger/internal/repo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBLoggerStorage struct {
	client *mongo.Client
}

func NewMongoDBLoggerStorage(client *mongo.Client) repo.ILoggerStorage {
	return &MongoDBLoggerStorage{
		client: client,
	}
}

func (l *MongoDBLoggerStorage) InsertOne(ctx context.Context, entry dto.LogEntry) error {
	collection := l.client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(ctx, dto.LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	return nil
}

func (r *MongoDBLoggerStorage) UpdateOne(ctx context.Context, entry dto.LogEntry) (MatchedCount int64, ModifiedCount int64, UpsertedCount int64, UpsertedID interface{}, err error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	collection := r.client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(entry.ID)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", entry.Name},
				{"data", entry.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		return 0, 0, 0, 0, err
	}

	return result.MatchedCount, result.ModifiedCount, result.UpsertedCount, result.UpsertedID, nil
}

func (l *MongoDBLoggerStorage) SelectAll(ctx context.Context) ([]*dto.LogEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	collection := l.client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*dto.LogEntry

	for cursor.Next(ctx) {
		var item dto.LogEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Print("Error decoding log into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	return logs, nil
}

func (l *MongoDBLoggerStorage) SelectOne(ctx context.Context, id string) (*dto.LogEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	collection := l.client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry dto.LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (l *MongoDBLoggerStorage) DropCollection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	collection := l.client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}
