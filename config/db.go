package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var Products *mongo.Collection

func init() {
	// get a mongo context
	clientOptions := options.Client().ApplyURI("mongodb+srv://sapling:nHYuR4vtCeKUWPfz@cluster0-qe9ag.gcp.mongodb.net/test?w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil { log.Fatal(err) }

	Products = client.Database("sapling").Collection("products")
	fmt.Println("You connected to your mongo database.")
}
