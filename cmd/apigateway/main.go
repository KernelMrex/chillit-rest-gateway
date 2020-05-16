package main

import (
	"chillit-rest-gateway/internal/app/apiserver"
	"chillit-rest-gateway/internal/app/configuration"
	"chillit-rest-gateway/internal/app/places"
	"flag"
	"log"

	"google.golang.org/grpc"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config_path", "config.yaml", "path for '.yaml' configuration file")
}

func main() {
	flag.Parse()

	log.Println("Parsing config")
	config, err := configuration.ParseConfig(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connecting places storage")
	conn, err := grpc.Dial(config.StoreService.URL, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	log.Println("Starting HTTP server")
	if err := apiserver.Start(config.APIServer, places.NewPlacesStoreClient(conn)); err != nil {
		log.Fatalln(err)
	}
}
