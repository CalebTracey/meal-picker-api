package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Cuisine struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Name   string             `bson:"name,omitempty" json:"name,omitempty"`
	Type   string             `bson:"type,omitempty" json:"type,omitempty"`
	Dishes []Dish             `bson:"dishes,omitempty" json:"dishes,omitempty"`
	Tags   []string           `bson:"tags,omitempty" json:"tags,omitempty"`
}

type Dish struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Cuisine primitive.ObjectID `bson:"cuisine,omitempty" json:"cuisine,omitempty"`
	Name    string             `bson:"name,omitempty" json:"name,omitempty"`
	Tags    []string           `bson:"tags,omitempty" json:"tags,omitempty"`
}
