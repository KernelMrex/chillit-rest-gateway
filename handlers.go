package main

import (
	"chillit-rest-gateway/models"
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

	// TODO: Request using gRPC

})

var NotImplementedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not implemented"))
	w.WriteHeader(http.StatusNotImplemented)
})
