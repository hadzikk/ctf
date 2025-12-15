package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SolvedChallenge struct {
	ChallengeID primitive.ObjectID `bson:"challenge" json:"challengeId"`
	SolvedAt    time.Time          `bson:"solvedAt" json:"solvedAt"`
	Points      int                `bson:"points" json:"points"`
}

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username         string             `bson:"username" json:"username" validate:"required,min=3,max=30"`
	Email            string             `bson:"email" json:"email" validate:"required,email"`
	Password         string             `bson:"password" json:"-" validate:"required,min=8"`
	Role             string             `bson:"role" json:"role" validate:"oneof=user admin"`
	Score            int                `bson:"score" json:"score"`
	SolvedChallenges []SolvedChallenge  `bson:"solvedChallenges" json:"solvedChallenges"`
	CreatedAt        time.Time          `bson:"createdAt" json:"createdAt"`
	LastActive       time.Time          `bson:"lastActive" json:"lastActive"`
}

func (u *User) BeforeCreate() {
	u.CreatedAt = time.Now()
	u.LastActive = time.Now()
	if u.Role == "" {
		u.Role = "user"
	}
	if u.Score == 0 {
		u.Score = 0
	}
}
