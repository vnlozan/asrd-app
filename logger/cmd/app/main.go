package main

import (
	"context"
	"log"
	"logger/internal/config"
	"logger/internal/controller"
	repo "logger/internal/repo/mongodb"
	"logger/internal/server"
	"logger/internal/service"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config := config.NewConfig()

	mongoClient, err := connectToMongo(config)
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	loggerStorage := repo.NewMongoDBLoggerStorage(mongoClient)
	loggerService := service.NewLoggerService(loggerStorage)
	loggerController := controller.NewLoggerController(loggerService)

	server := server.NewServer(config, loggerController)
	server.Start()
}

func connectToMongo(config *config.Config) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(config.MongoDBConfig.MongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return c, nil
}
