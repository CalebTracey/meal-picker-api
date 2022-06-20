package mongodb

import (
	"food-roulette-api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MapperI interface {
	MapDishesResponse(newIds []any, dishes []models.Dish, cuisineId primitive.ObjectID) []models.Dish
}

type Mapper struct{}

func (m Mapper) MapDishesResponse(newIds []any, dishes []models.Dish, cuisineId primitive.ObjectID) []models.Dish {
	response := make([]models.Dish, len(newIds))

	for i, v := range newIds {
		id := v.(primitive.ObjectID)
		response[i].ID = id
		response[i].Name = dishes[i].Name
		response[i].Cuisine = cuisineId
		response[i].Tags = dishes[i].Tags
	}

	return response
}
