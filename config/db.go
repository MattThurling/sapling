package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var CTX context.Context
var Products *mongo.Collection

func init() {
	// get a mongo context
	CTX, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(CTX, options.Client().ApplyURI("mongodb+srv://sapling:nHYuR4vtCeKUWPfz@cluster0-qe9ag.gcp.mongodb.net/test?w=majority"))
	if err != nil { log.Fatal(err) }

	Products = client.Database("sapling").Collection("products")
	fmt.Println("You connected to your mongo database.")
}
