package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AddCuisineRequest struct {
	Name   string   `json:"name,omitempty"`
	Dishes []Dish   `json:"dishes,omitempty"`
	Tags   []string `json:"tags,omitempty"`
}

type AddDishesRequest struct {
	Cuisine primitive.ObjectID `json:"cuisine,omitempty"`
	Name    string             `json:"name,omitempty"`
	Dishes  []Dish             `json:"dishes,omitempty"`
}
