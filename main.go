package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
)

type simpleCommentsAPI struct {
	topContext context.Context
	dbConn     *mongo.Client
	database   *mongo.Database

	// Actual models
	comments *mongo.Collection
}

func main() {
	topCtx := context.Background()

	// Create a new MongoDB client and connect to the server
	connURL := os.Getenv("SCA_MONGODB_URL")
	if len(connURL) == 0 {
		fmt.Println("No MongoDB URL found, exiting.")
		os.Exit(1)
	}

	client, err := mongo.Connect(topCtx, options.Client().ApplyURI(connURL))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(topCtx); err != nil {
			panic(err)
		}
	}()
	// Ping the primary
	if err := client.Ping(topCtx, readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to database, initializing.")

	db := client.Database("simple_comments")
	coll := db.Collection("comments")

	api := simpleCommentsAPI{
		topContext: topCtx,
		dbConn: client,
		database: db,
		comments: coll,
	}

	r := gin.Default()


	// Routes
	r.GET("/api/v1/:id/get-comments", api.getCommentsForURL)


	err = r.Run()
	if err != nil {
		panic(err)
	}
}
