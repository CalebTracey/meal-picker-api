package routes

import (
	"context"
	"encoding/json"
	"food-roulette-api/internal/facade"
	"food-roulette-api/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandler_HealthCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFacade := facade.NewMockServiceI(ctrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/add/cuisine", nil)

	type Test struct {
		name     string
		Service  facade.ServiceI
		want     http.HandlerFunc
		wantCode int
		wantRes  string
		r        *http.Request
		w        http.ResponseWriter
	}
	test := Test{
		name:     "Happy Path",
		Service:  mockFacade,
		wantCode: http.StatusOK,
		wantRes:  `{"ok": true}`,
		r:        r,
		w:        w,
	}
	t.Run(test.name, func(t *testing.T) {
		h := Handler{
			Service: test.Service,
		}

		test.w.Header().Set("Content-Type", "application/json")
		test.w.WriteHeader(test.wantCode)
		h.HealthCheck().ServeHTTP(test.w, test.r)
		res := w.Result()

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Errorf("error closing recorder body")
			}
		}(res.Body)
		resBody, _ := ioutil.ReadAll(res.Body)

		assert.Equal(t, test.wantCode, res.StatusCode)
		assert.JSONEq(t, test.wantRes, string(resBody))
	})
}

func TestHandler_AddNewCuisine_Happy_Path(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFacade := facade.NewMockServiceI(ctrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/add/cuisine", nil)

	type Test struct {
		name       string
		Service    facade.ServiceI
		wantCode   int
		ctx        context.Context
		wantRes    models.CuisineResponse
		cuisineReq models.AddCuisineRequest
		r          *http.Request
		w          http.ResponseWriter
	}
	test := Test{
		name:       "Happy Path",
		Service:    mockFacade,
		ctx:        context.Background(),
		wantCode:   http.StatusOK,
		r:          r,
		w:          w,
		cuisineReq: models.AddCuisineRequest{},
		wantRes: models.CuisineResponse{
			Cuisine: &models.Cuisine{},
		},
	}

	t.Run(test.name, func(t *testing.T) {
		h := Handler{
			Service: test.Service,
		}

		mockFacade.EXPECT().AddCuisine(test.ctx, test.cuisineReq).Return(test.wantRes).MaxTimes(1)
		test.w.Header().Set("Content-Type", "application/json")
		test.w.WriteHeader(test.wantCode)
		h.AddNewCuisine().ServeHTTP(test.w, test.r)
		res := w.Result()

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Errorf("error closing recorder body")
			}
		}(res.Body)

		var actualRes models.CuisineResponse

		err := json.NewDecoder(w.Body).Decode(&actualRes)
		if err != nil {
			t.Errorf("expected json to decode, got err: %v", err.Error())
		}

		assert.Equal(t, test.wantCode, res.StatusCode)
		assert.Equal(t, len(test.wantRes.Message.ErrorLog), len(actualRes.Message.ErrorLog))
	})

}

func TestHandler_AddNewCuisine_Sad_Path(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFacade := facade.NewMockServiceI(ctrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/add/cuisine", nil)

	type Test struct {
		name       string
		Service    facade.ServiceI
		wantCode   int
		ctx        context.Context
		wantRes    models.CuisineResponse
		cuisineReq models.AddCuisineRequest
		r          *http.Request
		w          http.ResponseWriter
	}
	test := Test{
		name:       "Sad Path",
		Service:    mockFacade,
		ctx:        context.Background(),
		wantCode:   http.StatusBadRequest,
		r:          r,
		w:          w,
		cuisineReq: models.AddCuisineRequest{},
		wantRes: models.CuisineResponse{
			Cuisine: nil,
			Message: models.Message{
				Status: strconv.Itoa(http.StatusBadRequest),
				ErrorLog: []models.ErrorLog{
					{
						Status:    strconv.Itoa(http.StatusBadRequest),
						RootCause: "Validation error",
						Trace:     "missing params for database insert",
					},
				},
			},
		},
	}

	t.Run(test.name, func(t *testing.T) {
		h := Handler{
			Service: test.Service,
		}
		mockFacade.EXPECT().AddCuisine(test.ctx, test.cuisineReq).Return(test.wantRes).MaxTimes(1)
		test.w.Header().Set("Content-Type", "application/json")
		test.w.WriteHeader(test.wantCode)
		h.AddNewCuisine().ServeHTTP(test.w, test.r)
		res := w.Result()

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Errorf("error closing recorder body")
			}
		}(res.Body)

		var actualRes models.CuisineResponse

		err := json.NewDecoder(w.Body).Decode(&actualRes)
		if err != nil {
			t.Errorf("expected json to decode, got err: %v", err.Error())
		}
		assert.Equal(t, test.wantCode, res.StatusCode)
		assert.Equal(t, len(test.wantRes.Message.ErrorLog), len(actualRes.Message.ErrorLog))
	})
}

func TestHandler_GetAllCuisines_Happy_Path(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFacade := facade.NewMockServiceI(ctrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/all/cuisines", nil)

	type Test struct {
		name     string
		Service  facade.ServiceI
		wantCode int
		ctx      context.Context
		wantRes  models.AllCuisinesResponse
		r        *http.Request
		w        http.ResponseWriter
	}
	test := Test{
		name:     "Happy Path",
		Service:  mockFacade,
		ctx:      context.Background(),
		wantCode: http.StatusOK,
		r:        r,
		w:        w,
		wantRes: models.AllCuisinesResponse{
			Cuisines: []*models.Cuisine{},
		},
	}
	t.Run(test.name, func(t *testing.T) {
		h := Handler{
			Service: test.Service,
		}
		mockFacade.EXPECT().AllCuisines(test.ctx).Return(test.wantRes).MaxTimes(1)
		test.w.Header().Set("Content-Type", "application/json")
		test.w.WriteHeader(test.wantCode)
		h.GetAllCuisines().ServeHTTP(test.w, test.r)
		res := w.Result()

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Errorf("error closing recorder body")
			}
		}(res.Body)

		var actualRes models.AllCuisinesResponse

		err := json.NewDecoder(w.Body).Decode(&actualRes)
		if err != nil {
			t.Errorf("expected json to decode, got err: %v", err.Error())
		}

		assert.Equal(t, test.wantCode, res.StatusCode)
		assert.Equal(t, len(test.wantRes.Message.ErrorLog), len(actualRes.Message.ErrorLog))
	})

}

func TestHandler_GetAllCuisines_Sad_Path(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFacade := facade.NewMockServiceI(ctrl)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/all/cuisines", nil)

	type Test struct {
		name     string
		Service  facade.ServiceI
		wantCode int
		ctx      context.Context
		wantRes  models.AllCuisinesResponse
		r        *http.Request
		w        http.ResponseWriter
	}
	test := Test{
		name:     "Sad Path",
		Service:  mockFacade,
		ctx:      context.Background(),
		wantCode: http.StatusInternalServerError,
		r:        r,
		w:        w,
		wantRes: models.AllCuisinesResponse{
			Message: models.Message{
				Status: strconv.Itoa(http.StatusBadRequest),
				ErrorLog: []models.ErrorLog{
					{
						Status:    strconv.Itoa(http.StatusBadRequest),
						RootCause: "FindAll error",
						Trace:     "missing params for database insert",
					},
				},
			},
		},
	}

	t.Run(test.name, func(t *testing.T) {
		h := Handler{
			Service: test.Service,
		}

		mockFacade.EXPECT().AllCuisines(test.ctx).Return(test.wantRes).MaxTimes(1)
		test.w.Header().Set("Content-Type", "application/json")
		test.w.WriteHeader(test.wantCode)
		h.GetAllCuisines().ServeHTTP(test.w, test.r)
		res := w.Result()

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Errorf("error closing recorder body")
			}
		}(res.Body)

		var actualRes models.AllCuisinesResponse

		err := json.NewDecoder(w.Body).Decode(&actualRes)
		if err != nil {
			t.Errorf("expected json to decode, got err: %v", err.Error())
		}
		assert.Equal(t, test.wantCode, res.StatusCode)
		assert.Equal(t, len(test.wantRes.Message.ErrorLog), len(actualRes.Message.ErrorLog))
	})
}
