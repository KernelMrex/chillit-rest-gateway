package apiserver

import (
	"chillit-rest-gateway/internal/app/places"
	"errors"
	"net/http"
)

// Start API web server
func Start(config *Config, placesStore places.PlacesStoreClient) error {
	srv := newServer(placesStore)
	if config == nil {
		return errors.New("apiserver could not start error: <nil> config")
	}
	return http.ListenAndServe(config.Hostname, srv)
}
