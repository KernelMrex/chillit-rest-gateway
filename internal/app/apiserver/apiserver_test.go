package apiserver

import (
	"bytes"
	"chillit-rest-gateway/internal/app/places"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	grpcServerPort = "10324"
)

func TestServer_postPlaceHandler_AllClear(t *testing.T) {
	// Building places client
	mockStoreClient := &places.MockPlacesStoreClient{
		AddPlaceResponse: &places.AddPlaceResponse{
			Id: 1,
		},
		AddPlaceError: nil,
	}

	// Building api server
	server := newServer(mockStoreClient)

	// Executing request on test server
	body := bytes.NewBuffer([]byte(`{
		"title": "test title",
		"address": "test address",
		"description": "test description"
	}`))
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/place", body)
	server.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestServer_postPlaceHandler_BadJSON(t *testing.T) {
	// Building places client
	mockStoreClient := &places.MockPlacesStoreClient{}

	// Building api server
	server := newServer(mockStoreClient)

	// Executing request on test server
	body := bytes.NewBuffer([]byte(`{bad json format}`))
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/place", body)
	server.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestServer_postPlaceHandler_UnavailablePlacesStore(t *testing.T) {
	// TODO: use grpc.clinettimeout error

	// Building places client
	mockStoreClient := &places.MockPlacesStoreClient{
		AddPlaceResponse: &places.AddPlaceResponse{},
		AddPlaceError:    errors.New("client timeout"),
	}

	// Building api server
	server := newServer(mockStoreClient)

	// Executing request on test server
	body := bytes.NewBuffer([]byte(`{
		"title": "test title",
		"address": "test address",
		"description": "test description"
	}`))
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/place", body)
	server.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadGateway, rec.Code)
}
