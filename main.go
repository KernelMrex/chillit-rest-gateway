package main

import (
	"chillit-rest-gateway/configuration"
	"chillit-rest-gateway/environment"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var Env *environment.Env

func init() {
	// Logger initialization
	initInfoLogger := log.New(os.Stdout, "Init info: ", 0)
	initErrorLogger := log.New(os.Stderr, "Init error: ", 0)
	initInfoLogger.Println("Initialization started...")

	// Getting config path from flag
	var confPath string
	flag.StringVar(&confPath, "config_path", "config.yaml", "path for '.yaml' configuration file")
	flag.Parse()

	// Build config and env
	conf, err := configuration.NewConfig(confPath)
	if err != nil {
		initErrorLogger.Fatalln(err)
	}
	initInfoLogger.Println("Configuration has loaded")

	Env, err = environment.BuildEnv(conf)
	if err != nil {
		initErrorLogger.Fatalln(err)
	}
	initInfoLogger.Println("Environment has built")
	initInfoLogger.Println("Initialization successful")
}

func main() {
	router := mux.NewRouter()

	router.Handle("/place", PostPlaceHandler).Methods(http.MethodPost)
	router.Handle("/places", NotImplementedHandler).Methods(http.MethodGet)

	if err := http.ListenAndServe(":8080", router); err != nil {
		Env.ErrorLogger.Fatalln("[ main ] error while serving:", err)
	}
}
