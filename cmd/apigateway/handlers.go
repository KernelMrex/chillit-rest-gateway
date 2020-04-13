package main

import (
	"chillit-rest-gateway/internal/app/models"
	"chillit-rest-gateway/internal/app/places"
	"encoding/json"
	"net/http"
)

var PostPlaceHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var place models.Place
	if err := decoder.Decode(&place); err != nil {
		Env.ErrorLogger.Println("[ PostPlaceHandler ] error while decoding request body:", err)
		return
	}

	_, err := Env.StoreService.AddPlace(r.Context(), &places.AddPlaceRequest{
		Place: &places.Place{
			Title:       place.Title,
			Address:     place.Address,
			Description: place.Description,
		},
	})
	if err != nil {
		Env.ErrorLogger.Println("[ PostPlaceHandler ] error while sending request to store service:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
})

var NotImplementedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not implemented"))
	return
})
