package main

import (
	_ "embed"
	"log"
	"net/http"
	"time"

	"github.com/Bronku/iroon/server"
	"github.com/BurntSushi/toml"
)

type config struct {
	Database struct {
		File string
	}

	Server struct {
		Addr              string
		ReadHeaderTimeout time.Duration
	}
}

//go:embed config.default.toml
var defaultConfig string

func main() {
	// load config
	var conf config
	_, err := toml.Decode(defaultConfig, &conf)
	if err != nil {
		log.Fatal(err)
	}
	_, err = toml.DecodeFile("config.toml", &conf)
	if err != nil {
		log.Println("using default config")
	}

	// load database
	db, err := store.loadStore(conf.Database.File)
	if err != nil {
		log.Panic(err)
	}

	// load server
	h := server.New(db)
	var handler http.Handler = h
	log.Println("starting server")

	// start server
	server := &http.Server{
		Handler:           handler,
		Addr:              conf.Server.Addr,
		ReadHeaderTimeout: conf.Server.ReadHeaderTimeout,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
