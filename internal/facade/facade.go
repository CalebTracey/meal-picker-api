package facade

import (
	"context"
	"fmt"
	"food-roulette-api/internal/models"
	"food-roulette-api/internal/services/mongodb"
	config "github.com/calebtracey/config-yaml"
	"net/http"
	"strconv"
)

//go:generate mockgen -destination=mockFacade.go -package=facade . ServiceI
type ServiceI interface {
	AddCuisine(ctx context.Context, cuisine models.AddCuisineRequest) models.CuisineResponse
	AllCuisines(ctx context.Context) models.AllCuisinesResponse
}

type Service struct {
	MongoService mongodb.ServiceI
}

func NewService(appConfig *config.Config) (Service, error) {
	mongoService, err := mongodb.InitializeMongoService(appConfig)
	if err != nil {
		return Service{}, err
	}
	return Service{
		MongoService: mongoService,
	}, nil
}

func (s *Service) AddCuisine(ctx context.Context, cuisine models.AddCuisineRequest) (response models.CuisineResponse) {
	var message models.Message

	if cuisine.Name == "" {
		message.ErrorLog = errorLogs([]error{fmt.Errorf("missing params for database insert")}, "Validation error", http.StatusBadRequest)
		message.Status = strconv.Itoa(http.StatusBadRequest)
		response.Message = message
		return response
	}

	result, err := s.MongoService.AddNewCuisine(ctx, cuisine)
	if err != nil {
		message.ErrorLog = errorLogs([]error{err}, "Insertion error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
		response.Message = message
		return response
	}

	response.Cuisine = result
	response.Message.Status = strconv.Itoa(http.StatusOK)

	return response
}

func (s *Service) AllCuisines(ctx context.Context) (response models.AllCuisinesResponse) {
	var message models.Message

	results, err := s.MongoService.GetAllCuisines(ctx)
	if err != nil {
		message.ErrorLog = errorLogs([]error{err}, "FindAll error", http.StatusInternalServerError)
		message.Status = strconv.Itoa(http.StatusInternalServerError)
		response.Message = message
		return response
	}
	response.Message.Status = strconv.Itoa(http.StatusOK)
	response.Cuisines = results

	return response
}

func errorLogs(errors []error, rootCause string, status int) []models.ErrorLog {
	var errLogs []models.ErrorLog
	for _, err := range errors {
		errLogs = append(errLogs, models.ErrorLog{
			RootCause: rootCause,
			Status:    strconv.Itoa(status),
			Trace:     err.Error(),
		})
	}
	return errLogs
}
