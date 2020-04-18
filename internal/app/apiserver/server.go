package apiserver

import (
	"chillit-rest-gateway/internal/app/places"
	"encoding/json"
	"net/http"
	"strconv"

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
	s.router.HandleFunc("/places", s.getPlacesHandler()).Methods(http.MethodGet)
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
			w.WriteHeader(http.StatusBadRequest)
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
			w.WriteHeader(http.StatusBadGateway)
			return
		}
	})
}

func (s *server) getPlacesHandler() http.HandlerFunc {
	const methodName string = "GET /places"

	// Got as URLQuery
	type request struct {
		Offset uint64
		Amount uint64
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
		var reqPlaces request
		urlVals := r.URL.Query()

		var err error
		reqPlaces.Offset, err = strconv.ParseUint(urlVals.Get("offset"), 10, 64)
		if err != nil {
			s.logger.Errorf("bad parameters for '%s' request: amount='%v'", methodName, r.URL.RawQuery)
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		reqPlaces.Amount, err = strconv.ParseUint(urlVals.Get("amount"), 10, 64)
		if err != nil || reqPlaces.Amount == 0 {
			s.logger.Errorf("bad parameters for '%s' request: amount='%v'", methodName, r.URL.RawQuery)
			w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		placesStoreResp, err := s.placesStore.GetPlaces(r.Context(), &places.GetPlacesRequest{
			Amount: reqPlaces.Amount,
			Offset: reqPlaces.Offset,
		})
		if err != nil {
			s.logger.Errorf("error in '%s' while requesting places from placesStore: %v", methodName, err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		respPlaces := make([]*responsePlace, len(placesStoreResp.Places))
		for i, place := range placesStoreResp.Places {
			respPlaces[i] = &responsePlace{
				ID:          place.Id,
				Title:       place.Title,
				Address:     place.Address,
				Description: place.Description,
			}
		}

		if err := json.NewEncoder(w).Encode(response{Places: respPlaces}); err != nil {
			s.logger.Errorf("error in '%s' while encoding json response", methodName)
			w.WriteHeader(http.StatusInternalServerError)
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
