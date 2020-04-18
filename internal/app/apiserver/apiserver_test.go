package apiserver

import (
	"bytes"
	"chillit-rest-gateway/internal/app/places"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	grpcServerPort = "10324"
)

func TestServer_postPlaceHandler_Success(t *testing.T) {
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

func TestServer_getPlacesHandler_Success(t *testing.T) {
	// Building places client
	mockStoreClient := &places.MockPlacesStoreClient{
		GetPlacesResponse: &places.GetPlacesResponse{
			Places: []*places.Place{
				{
					Id:          1,
					Title:       "test title 1",
					Address:     "test address 1",
					Description: "test description 1",
				},
				{
					Id:          2,
					Title:       "test title 2",
					Address:     "test address 2",
					Description: "test description 2",
				},
			},
		},
	}

	// Building api server
	server := newServer(mockStoreClient)

	// Executing request on test server
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/places?offset=0&amount=20", nil)
	server.ServeHTTP(rec, req)

	// Checking response
	type responsePlace struct {
		ID          uint64 `json:"id"`
		Title       string `json:"title"`
		Address     string `json:"address"`
		Description string `json:"description"`
	}
	type response struct {
		Places []*responsePlace `json:"places"`
	}
	respBody, _ := json.Marshal(response{
		Places: []*responsePlace{
			{
				ID:          1,
				Title:       "test title 1",
				Address:     "test address 1",
				Description: "test description 1",
			},
			{
				ID:          2,
				Title:       "test title 2",
				Address:     "test address 2",
				Description: "test description 2",
			},
		},
	})

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, string(rec.Body.Bytes()), string(respBody)+"\n")
}
