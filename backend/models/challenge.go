package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MapPoint struct {
	Lat float64 `bson:"lat" json:"lat"`
	Lng float64 `bson:"lng" json:"lng"`
}

type MapMarker struct {
	Position    MapPoint `bson:"position" json:"position"`
	Title       string   `bson:"title" json:"title"`
	Description string   `bson:"description" json:"description"`
}

type MapConfig struct {
	Center  MapPoint    `bson:"center" json:"center"`
	Zoom    int         `bson:"zoom" json:"zoom"`
	Markers []MapMarker `bson:"markers" json:"markers"`
}

type File struct {
	Name     string `bson:"name" json:"name"`
	URL      string `bson:"url" json:"url"`
	Size     int64  `bson:"size" json:"size"`
	MimeType string `bson:"mimeType" json:"mimeType"`
}

type Hint struct {
	Text          string `bson:"text" json:"text"`
	PointsPenalty int    `bson:"pointsPenalty" json:"pointsPenalty"`
}

type Challenge struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title       string             `bson:"title" json:"title" validate:"required"`
	Description string             `bson:"description" json:"description" validate:"required"`
	Category    string             `bson:"category" json:"category" validate:"required,oneof=Web Cryptography Forensics 'Reverse Engineering' PWN Misc GIS"`
	Difficulty  string             `bson:"difficulty" json:"difficulty" validate:"required,oneof=Easy Medium Hard Expert"`
	Points      int                `bson:"points" json:"points" validate:"required,min=0"`
	Flag        string             `bson:"flag" json:"-" validate:"required"`
	Hints       []Hint             `bson:"hints,omitempty" json:"hints,omitempty"`
	MapConfig   *MapConfig         `bson:"mapConfig,omitempty" json:"mapConfig,omitempty"`
	Files       []File             `bson:"files,omitempty" json:"files,omitempty"`
	IsActive    bool               `bson:"isActive" json:"isActive"`
	AuthorID    primitive.ObjectID `bson:"author" json:"authorId" validate:"required"`
	Solves      int                `bson:"solves" json:"solves"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func (c *Challenge) BeforeCreate() {
	timeNow := time.Now()
	c.CreatedAt = timeNow
	c.UpdatedAt = timeNow
	if c.IsActive == false {
		c.IsActive = true
	}
	if c.Solves == 0 {
		c.Solves = 0
	}
}
