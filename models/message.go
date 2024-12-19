package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	FromUser  string             `bson:"fromUser"`
	ToUser    string             `bson:"toUser"`
	Content   string             `bson:"content"`
	Timestamp time.Time          `bson:"timestamp"`
}
