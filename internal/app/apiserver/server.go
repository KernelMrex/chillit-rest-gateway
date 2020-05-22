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
	logger         *logrus.Logger
	placesStore    places.PlacesStoreClient
	router         *mux.Router
	allowedOrigins string
}

func newServer(placesStore places.PlacesStoreClient, allowedOrigins string) *server {
	s := &server{
		logger:         logrus.New(),
		router:         mux.NewRouter(),
		placesStore:    placesStore,
		allowedOrigins: allowedOrigins,
	}
	s.configureRouter()
	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/places", s.CorsMiddleware(s.getPlacesHandler())).Methods(http.MethodGet)
	s.router.HandleFunc("/cities", s.CorsMiddleware(s.getCitiesHandler())).Methods(http.MethodGet)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// TODO: separate origin url in a config
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin")
		w.Header().Add("Access-Control-Allow-Origin", s.allowedOrigins)
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET")
		next.ServeHTTP(w, r)
	})
}

func (s *server) getPlacesHandler() http.HandlerFunc {
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
		ImgURL      string `json:"image_url"`
	}

	type response struct {
		Places []*responsePlace `json:"places"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestValues request
		queryDecoder := schema.NewDecoder()
		queryDecoder.IgnoreUnknownKeys(true)
		if err := queryDecoder.Decode(&requestValues, r.URL.Query()); err != nil {
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
		for i, pbPlace := range placesStoreResp.Places {
			jsonFormattableResponse.Places[i] = &responsePlace{
				ID:          pbPlace.GetId(),
				Title:       pbPlace.GetTitle(),
				Address:     pbPlace.GetAddress(),
				Description: pbPlace.GetDescription(),
				ImgURL:      pbPlace.GetImgURL(),
			}
		}

		if err := json.NewEncoder(w).Encode(&jsonFormattableResponse); err != nil {
			s.logger.Errorf("could not encode response, error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func (s *server) getCitiesHandler() http.HandlerFunc {
	type request struct {
		Offset uint64 `schema:"offset"`
		Amount uint64 `schema:"amount"`
	}

	type responseCity struct {
		ID    uint64 `json:"id"`
		Title string `json:"title"`
	}

	type response struct {
		Cities []*responseCity `json:"cities"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestValues request
		queryDecoder := schema.NewDecoder()
		queryDecoder.IgnoreUnknownKeys(true)
		if err := queryDecoder.Decode(&requestValues, r.URL.Query()); err != nil {
			s.logger.Errorf("could not decode GET query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		citiesStoreResp, err := s.placesStore.GetCities(context.Background(), &places.GetCitiesRequest{
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
			Cities: make([]*responseCity, len(citiesStoreResp.Cities)),
		}
		for i, pbCity := range citiesStoreResp.Cities {
			jsonFormattableResponse.Cities[i] = &responseCity{
				ID:    pbCity.GetId(),
				Title: pbCity.GetTitle(),
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
