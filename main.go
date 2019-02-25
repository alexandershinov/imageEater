package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Config struct {
	Port       int
	FilesField string
	FilesDir   string
}

var config Config

// read config from file (toml format)
func ReadConfig(configFile string) Config {
	defaultConfig := Config{
		Port:       4000,
		FilesField: "files",
		FilesDir:   "files/",
	}

	_, err := os.Stat(configFile)
	if err != nil {
		log.Printf("config file '%s' is missing\n", configFile)
		return defaultConfig
	}

	var config Config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Printf("config file '%s' parsing error, default configuration uses\n", configFile)
		return defaultConfig
	}

	return config
}

func listenAndServe(port string, router *mux.Router) *http.Server {
	srv := &http.Server{Addr: port, Handler: router}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()
	return srv
}

func main() {
	// try to read config
	configFile := flag.String("config", "config.toml", "config file path")
	flag.Parse()
	config = ReadConfig(*configFile)

	// add router
	router := mux.NewRouter()
	// add REST API method
	router.HandleFunc("/images", LoadImages).Methods("POST")

	// start http server
	log.Printf("starting HTTP server")
	srv := listenAndServe(fmt.Sprintf(":%d", config.Port), router)

	// Setting up signal capturing
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (pkill -2)
	<-stop

	//log.Printf("serving for 5 seconds")
	//time.Sleep(5 * time.Second)
	log.Printf("stopping server")
	// graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Printf("exit")
}
