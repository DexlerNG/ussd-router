package startups

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"sync"
)

var connectOnce sync.Once
var client *mongo.Client


func GetMongoClient() *mongo.Client {
	println("Calling .Once ......", os.Getenv("MONGODB_URL") != "")
	connectOnce.Do(func() {
		println("Connecting to Mongodb......")
		data, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
		if err != nil {
			panic(err)
		}
		err = data.Ping(context.TODO(), readpref.Primary())
		if err != nil {
			log.Fatalf("Exit error Ping: %v", err)
		}
		println("Connected to mongo....")
		client  = data
	})
	println("client", client)
	return client
}

func GetMongoDatabase() *mongo.Database {
	return client.Database(os.Getenv("DATABASE_NAME"))
}