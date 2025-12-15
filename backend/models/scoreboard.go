package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Scoreboard struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user" json:"userId" validate:"required"`
	Username  string             `bson:"username" json:"username" validate:"required"`
	Score     int                `bson:"score" json:"score"`
	LastSolve *time.Time         `bson:"lastSolve,omitempty" json:"lastSolve,omitempty"`
	Rank      int                `bson:"rank,omitempty" json:"rank,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
