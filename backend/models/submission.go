package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Submission struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID        primitive.ObjectID `bson:"user" json:"userId" validate:"required"`
	ChallengeID   primitive.ObjectID `bson:"challenge" json:"challengeId" validate:"required"`
	Flag          string             `bson:"flag" json:"flag" validate:"required"`
	IsCorrect     bool               `bson:"isCorrect" json:"isCorrect"`
	PointsAwarded int                `bson:"pointsAwarded" json:"pointsAwarded"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
}

func (s *Submission) BeforeCreate() {
	s.CreatedAt = time.Now()
}
