package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service struct {
	ID   primitive.ObjectID   `json:"id,omitempty" bson:"id,omitempty"`
	Name string               `json:"name" bson:"name"`
	Sub  []primitive.ObjectID `json:"sub" bson:"sub"`
}

type SubService struct {
	ID   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer func() {
		client.Disconnect(ctx)
		cancel()
	}()
	collection := client.Database("testing").Collection("service")
	sub1Result, err := collection.InsertOne(ctx, &SubService{
		Name: "sub1",
	})
	if err != nil {
		panic(err)
	}
	sub1ID, _ := sub1Result.InsertedID.(primitive.ObjectID)
	fmt.Println("Sub1 ID", sub1ID)
	sub2Result, err := collection.InsertOne(ctx, &SubService{
		Name: "sub2",
	})
	if err != nil {
		panic(err)
	}
	sub2ID, _ := sub2Result.InsertedID.(primitive.ObjectID)
	fmt.Println("Sub2 ID", sub2ID)
	serviceResult, err := collection.InsertOne(ctx, &Service{
		Name: "service",
		Sub:  []primitive.ObjectID{sub1ID, sub2ID},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Service ID", serviceResult.InsertedID)
	// Get service
	var s Service
	if err = collection.FindOne(ctx, bson.M{"_id": serviceResult.InsertedID}).Decode(&s); err != nil {
		panic(err)
	}
	fmt.Println("Service's sub", s.Sub)
}
