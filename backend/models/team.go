package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Team struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string               `bson:"name" json:"name" validate:"required,min=3"`
	Members   []primitive.ObjectID `bson:"members" json:"members"`
	Score     int                  `bson:"score" json:"score"`
	CaptainID primitive.ObjectID   `bson:"captain" json:"captainId" validate:"required"`
	Token     string               `bson:"token,omitempty" json:"token,omitempty"`
	CreatedAt time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt" json:"updatedAt"`
}

func (t *Team) BeforeCreate() {
	timeNow := time.Now()
	t.CreatedAt = timeNow
	t.UpdatedAt = timeNow
	if t.Score == 0 {
		t.Score = 0
	}
}
