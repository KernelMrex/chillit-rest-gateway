package apiserver

import (
	"chillit-rest-gateway/internal/app/places"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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

// TODO: add logging middleware
func (s *server) configureRouter() {
	s.router.HandleFunc("/place", s.postPlaceHandler()).Methods(http.MethodPost)
	s.router.HandleFunc("/places", s.notImplementedHandler()).Methods(http.MethodGet)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) postPlaceHandler() http.HandlerFunc {
	type requestPlace struct {
		Title       string `json:"title"`
		Address     string `json:"address"`
		Description string `json:"description"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var place requestPlace
		if err := decoder.Decode(&place); err != nil {
			s.logger.Errorln("[ PostPlaceHandler ] error while decoding request body:", err)
			return
		}
		_, err := s.placesStore.AddPlace(r.Context(), &places.AddPlaceRequest{
			Place: &places.Place{
				Title:       place.Title,
				Address:     place.Address,
				Description: place.Description,
			},
		})
		if err != nil {
			s.logger.Errorln("[ PostPlaceHandler ] error while sending request to store service:", err)
			w.WriteHeader(http.StatusBadRequest)
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
