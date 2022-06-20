package routes

import (
	"encoding/json"
	"food-roulette-api/internal/facade"
	"food-roulette-api/internal/models"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Handler struct {
	Service facade.ServiceI
}

func (h Handler) InitializeRoutes() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	// Health check
	r.Handle("/api/health", h.HealthCheck()).Methods(http.MethodGet)

	r.Handle("/api/all/cuisines", h.GetAllCuisines()).Methods(http.MethodGet)

	r.Handle("/api/add/cuisine", h.AddNewCuisine()).Methods(http.MethodPost)

	r.Handle("/api/add/all/dishes", h.AddDishes()).Methods(http.MethodPost)
	return r
}

func (h Handler) HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		if err != nil {
			logrus.Errorln(err.Error())
			return
		}
	}
}

func (h Handler) AddNewCuisine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		var response models.CuisineResponse

		defer func() {
			response, status := setInsertResponse(response)
			response.Message.TimeTaken = time.Since(startTime).String()
			_ = json.NewEncoder(writeHeader(w, status)).Encode(response)
		}()

		apiRequest := models.AddCuisineRequest{}
		requestBody, readErr := ioutil.ReadAll(r.Body)

		if readErr != nil {
			response.Message.ErrorLog = errorLogs([]error{readErr}, "Unable to read request body", http.StatusBadRequest)
			return
		}
		err := json.Unmarshal(requestBody, &apiRequest)
		if err != nil {
			response.Message.ErrorLog = errorLogs([]error{err}, "Unable to parse request", http.StatusBadRequest)

		}

		response = h.Service.AddCuisine(r.Context(), apiRequest)
	}
}

func (h Handler) AddDishes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (h Handler) GetAllCuisines() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		var response models.AllCuisinesResponse

		defer func() {
			response, status := setAllResponse(response)
			response.Message.TimeTaken = time.Since(startTime).String()
			_ = json.NewEncoder(writeHeader(w, status)).Encode(response)
		}()

		response = h.Service.AllCuisines(r.Context())
	}
}

func setAllResponse(res models.AllCuisinesResponse) (models.AllCuisinesResponse, int) {
	hn, _ := os.Hostname()
	status, _ := strconv.Atoi(res.Message.Status)
	res.Message.HostName = hn
	return res, status
}

func setInsertResponse(res models.CuisineResponse) (models.CuisineResponse, int) {
	hn, _ := os.Hostname()
	status, _ := strconv.Atoi(res.Message.Status)
	res.Message.HostName = hn
	return res, status
}

func writeHeader(w http.ResponseWriter, code int) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return w
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
