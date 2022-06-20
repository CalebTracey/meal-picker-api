package facade

import (
	"context"
	"fmt"
	"food-roulette-api/internal/models"
	"food-roulette-api/internal/services/mongodb"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reflect"
	"strconv"
	"testing"
)

func TestService_AddCuisine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMongoSvc := mongodb.NewMockServiceI(ctrl)
	happyCuisineDb := &models.Cuisine{
		ID:   primitive.ObjectID{},
		Name: "test food",
		Type: "cuisine",
	}

	tests := []struct {
		name          string
		MongoService  mongodb.ServiceI
		ctx           context.Context
		cuisine       models.AddCuisineRequest
		mockCuisineDb *models.Cuisine
		wantResponse  models.CuisineResponse
		wantError     error
	}{
		{
			name:         "Happy Path",
			MongoService: mockMongoSvc,
			ctx:          context.Background(),
			cuisine: models.AddCuisineRequest{
				Name: "test food",
			},
			mockCuisineDb: happyCuisineDb,
			wantResponse: models.CuisineResponse{
				Cuisine: happyCuisineDb,
				Message: models.Message{
					Status: strconv.Itoa(http.StatusOK),
				},
			},
			wantError: nil,
		},
		{
			name:         "Sad Path: missing param",
			MongoService: mockMongoSvc,
			ctx:          context.Background(),
			cuisine: models.AddCuisineRequest{
				Name: "",
			},
			wantResponse: models.CuisineResponse{
				Cuisine: nil,
				Message: models.Message{
					ErrorLog: []models.ErrorLog{
						{
							Status:    strconv.Itoa(http.StatusBadRequest),
							RootCause: "Validation error",
							Trace:     "missing params for database insert",
						},
					},
					Status: strconv.Itoa(http.StatusBadRequest),
				},
			},
			wantError: fmt.Errorf("missing params for database insert"),
		},
		{
			name:         "Sad Path: insert error",
			MongoService: mockMongoSvc,
			ctx:          context.Background(),
			cuisine: models.AddCuisineRequest{
				Name: "test food",
			},
			wantResponse: models.CuisineResponse{
				Cuisine: nil,
				Message: models.Message{
					ErrorLog: []models.ErrorLog{
						{
							Status:    strconv.Itoa(http.StatusInternalServerError),
							RootCause: "Insertion error",
							Trace:     "test error",
						},
					},
					Status: strconv.Itoa(http.StatusInternalServerError),
				},
			},
			wantError: fmt.Errorf("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				MongoService: tt.MongoService,
			}
			mockMongoSvc.EXPECT().AddNewCuisine(tt.ctx, tt.cuisine).Return(tt.mockCuisineDb, tt.wantError).MaxTimes(1)
			if gotResponse := s.AddCuisine(tt.ctx, tt.cuisine); !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("AddCuisine() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}

func TestService_AllCuisines(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMongoSvc := mongodb.NewMockServiceI(ctrl)
	happyCuisines := []*models.Cuisine{
		{
			primitive.ObjectID{},
			"test food one",
			"cuisine",
			[]models.Dish{},
			[]string{},
		},
		{
			primitive.ObjectID{},
			"test food two",
			"cuisine",
			[]models.Dish{},
			[]string{},
		},
	}

	tests := []struct {
		name            string
		MongoService    mongodb.ServiceI
		ctx             context.Context
		wantSvcResponse []*models.Cuisine
		wantResponse    models.AllCuisinesResponse
		wantError       error
	}{
		{
			name:            "Happy Path",
			MongoService:    mockMongoSvc,
			ctx:             context.Background(),
			wantSvcResponse: happyCuisines,
			wantResponse: models.AllCuisinesResponse{
				Cuisines: happyCuisines,
				Message: models.Message{
					Status: strconv.Itoa(http.StatusOK),
				},
			},
			wantError: nil,
		},
		{
			name:            "Sad Path: service error",
			MongoService:    mockMongoSvc,
			ctx:             context.Background(),
			wantSvcResponse: nil,
			wantResponse: models.AllCuisinesResponse{
				Message: models.Message{
					ErrorLog: []models.ErrorLog{
						{
							Status:    strconv.Itoa(http.StatusInternalServerError),
							RootCause: "FindAll error",
							Trace:     "test error",
						},
					},
					Status: strconv.Itoa(http.StatusInternalServerError),
				},
			},
			wantError: fmt.Errorf("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				MongoService: tt.MongoService,
			}

			mockMongoSvc.EXPECT().GetAllCuisines(tt.ctx).Return(tt.wantSvcResponse, tt.wantError).MaxTimes(1)
			if gotResponse := s.AllCuisines(tt.ctx); !reflect.DeepEqual(gotResponse, tt.wantResponse) {
				t.Errorf("AllCuisines() = %v, want %v", gotResponse, tt.wantResponse)
			}
		})
	}
}
