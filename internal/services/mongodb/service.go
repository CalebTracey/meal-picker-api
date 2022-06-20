package mongodb

import (
	"context"
	"fmt"
	"food-roulette-api/internal/models"
	config "github.com/calebtracey/config-yaml"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate mockgen -destination=mockService.go -package=mongodb . ServiceI
type ServiceI interface {
	AddNewCuisine(ctx context.Context, request models.AddCuisineRequest) (*models.Cuisine, error)
	AddAllDishes(ctx context.Context, request models.AddDishesRequest) ([]models.Dish, error)
	GetAllCuisines(ctx context.Context) ([]*models.Cuisine, error)
}

type Service struct {
	Database string
	Client   *mongo.Client
	Mapper   Mapper
}

func InitializeMongoService(appConfig *config.Config) (*Service, error) {
	mongoConfig, err := appConfig.GetDatabaseConfig("MONGO")
	if err != nil {
		return nil, err
	}
	return &Service{
		Database: mongoConfig.Database.Value,
		Client:   mongoConfig.MongoClient,
	}, nil

}

func (s *Service) AddNewCuisine(ctx context.Context, request models.AddCuisineRequest) (*models.Cuisine, error) {
	dbName := s.Database
	database := s.Client.Database(dbName)
	cuisineColl := database.Collection("cuisines")
	var response models.Cuisine
	var dishes []models.Dish
	var cuisineId any
	var err error

	found, err := cuisineColl.Find(ctx, bson.M{"name": request.Name})
	if err != nil {
		return &response, err
	}

	if found.RemainingBatchLength() > 0 {
		return &response, fmt.Errorf("%v already exists in the database", request.Name)
	}

	cursor, err := database.Collection("cuisines").InsertOne(ctx, request)
	if err != nil {
		return &response, err
	}
	log.Infof("inserted new cuisine: %v into database", request.Name)
	if cursor.InsertedID != nil {
		cuisineId = cursor.InsertedID
	}

	//TODO refactor this into UpdateDishes func
	if len(request.Dishes) > 0 {
		dishRequest := models.AddDishesRequest{
			Cuisine: cuisineId.(primitive.ObjectID),
			Dishes:  request.Dishes,
		}
		dishes, err = s.AddAllDishes(ctx, dishRequest)
		if err != nil {
			return &response, err
		}
		request.Dishes = dishes
		update := bson.M{
			"$set": request,
		}
		cursor, collErr := database.Collection("cuisines").UpdateByID(ctx, cursor.InsertedID, update)
		if collErr != nil {
			return &response, collErr
		}
		log.Infof("update new cuisine: %v with new dishes", request.Name)
		if cursor.UpsertedID != nil {
			cuisineId = cursor.UpsertedID
		}
	}

	response = models.Cuisine{
		ID:     cuisineId.(primitive.ObjectID),
		Name:   request.Name,
		Dishes: request.Dishes,
		Tags:   request.Tags,
	}

	return &response, nil
}

func (s *Service) AddAllDishes(ctx context.Context, request models.AddDishesRequest) ([]models.Dish, error) {
	dbName := s.Database
	database := s.Client.Database(dbName)
	var docs []interface{}
	var results []models.Dish
	var err error

	for _, dish := range request.Dishes {
		doc, docErr := toDoc(dish)
		if docErr != nil {
			return nil, docErr
		}
		docs = append(docs, doc)
	}

	cursor, err := database.Collection("dishes").InsertMany(ctx, docs)
	if err != nil {
		return results, err
	}

	results = s.Mapper.MapDishesResponse(cursor.InsertedIDs, request.Dishes, request.Cuisine)

	return results, nil
}

func (s *Service) GetAllCuisines(ctx context.Context) ([]*models.Cuisine, error) {
	dbName := s.Database
	database := s.Client.Database(dbName)
	var results []*models.Cuisine
	var err error

	cursor, err := database.Collection("cuisines").Find(ctx, bson.D{})
	if err != nil {
		return results, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			err = fmt.Errorf("failed to close mongodb cursor; err: %v", err.Error())
		}
	}(cursor, ctx)

	if curErr := cursor.All(ctx, &results); curErr != nil {
		return nil, curErr
	}

	return results, nil
}

func toDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}
