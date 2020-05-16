package apiserver

import (
	"chillit-rest-gateway/internal/app/places"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/sirupsen/logrus"
)

type server struct {
	logger      *logrus.Logger
	placesStore places.PlacesStoreClient
	router      *mux.Router
}

func newServer(placesStore places.PlacesStoreClient) *server {
	s := &server{
		logger:      logrus.New(),
		router:      mux.NewRouter(),
		placesStore: placesStore,
	}
	s.configureRouter()
	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/places", s.getPlacesHandler()).Methods(http.MethodGet)
	s.router.HandleFunc("/cities", s.notImplementedHandler()).Methods(http.MethodGet)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) getPlacesHandler() http.HandlerFunc {
	const methodName string = "GET /places"

	type request struct {
		Offset uint64 `schema:"offset"`
		Amount uint64 `schema:"amount"`
		CityID uint64 `schema:"city_id"`
	}

	type responsePlace struct {
		ID          uint64 `json:"id"`
		Title       string `json:"title"`
		Address     string `json:"address"`
		Description string `json:"description"`
	}

	type response struct {
		Places []*responsePlace `json:"places"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestValues request
		if err := schema.NewDecoder().Decode(&requestValues, r.URL.Query()); err != nil {
			s.logger.Errorf("could not decode GET query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Get places by city name
		placesStoreResp, err := s.placesStore.GetPlacesByCityID(context.Background(), &places.GetPlacesByCityIDRequest{
			CityID: requestValues.CityID,
			Amount: requestValues.Amount,
			Offset: requestValues.Offset,
		})
		if err != nil {
			s.logger.Errorf("could not get data from places store, error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Converting PB to JSON
		jsonFormattableResponse := response{
			Places: make([]*responsePlace, len(placesStoreResp.Places)),
		}
		for i, grpcPlace := range placesStoreResp.Places {
			jsonFormattableResponse.Places[i] = &responsePlace{
				ID:          grpcPlace.GetId(),
				Title:       grpcPlace.GetTitle(),
				Address:     grpcPlace.GetAddress(),
				Description: grpcPlace.GetDescription(),
			}
		}

		if err := json.NewEncoder(w).Encode(&jsonFormattableResponse); err != nil {
			s.logger.Errorf("could not encode response, error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *server) notImplementedHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Not implemented"))
		return
	})
}
