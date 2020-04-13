package environment

import (
	"chillit-rest-gateway/internal/app/configuration"
	"chillit-rest-gateway/internal/app/places"
	"errors"
	"google.golang.org/grpc"
	"log"
	"os"
)

type Env struct {
	InfoLogger   *log.Logger
	ErrorLogger  *log.Logger
	StoreService places.PlacesStoreClient
}

func BuildEnv(config *configuration.Configuration) (*Env, error) {
	storeService, err := connectStoreService(config.StoreServiceConfig.Url)
	if err != nil {
		return nil, err
	}

	return &Env{
		InfoLogger:   log.New(os.Stdout, "INFO: ", 0),
		ErrorLogger:  log.New(os.Stderr, "ERROR:", 0),
		StoreService: storeService,
	}, nil
}

func connectStoreService(url string) (places.PlacesStoreClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, errors.New("[ connectStoreService ] could not connect: " + err.Error())
	}
	return places.NewPlacesStoreClient(conn), nil
}
