package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type URL struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty"`
	ShortCode     string              `json:"short_code" bson:"short_code" unique:"true"`
	ShortURL      string              `json:"short_url" bson:"short_url"`
	LongURL       string              `json:"long_url" bson:"long_url"`
	CreatedAt     time.Time           `json:"created_at" bson:"created_at"`
	ExpiresAt     *time.Time          `json:"expires_at" bson:"expires_at,omitempty"`
	UserID        *primitive.ObjectID `json:"user_id" bson:"user_id,omitempty"`
	Clicks        int64               `json:"clicks" bson:"clicks"`
	LastClickedAt *time.Time          `json:"last_clicked_at" bson:"last_clicked_at,omitempty"`
}

var URLCollection *mongo.Collection

func SetupURLCollection(client *mongo.Client) {
	URLCollection = client.Database("goapp").Collection("urls")

	// Create a unique index on short_code
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "short_code", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := URLCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatalf("Failed to create unique index for short_code: %v", err)
	}
}
