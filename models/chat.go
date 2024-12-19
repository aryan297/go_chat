package models

import "time"

type Chat struct {
	ID        string    `bson:"_id,omitempty"`
	UserID    string    `bson:"userId"`
	Message   string    `bson:"message"`
	Timestamp time.Time `bson:"timestamp"`
}
