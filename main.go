package main

import (
	"github.com/dylan-dinh/api-gin/pkg/conf"
	"github.com/dylan-dinh/api-gin/pkg/database"
	"github.com/dylan-dinh/api-gin/pkg/handlers"
	"log"
	"net/http"
)

func main() {
	db, err := database.NewDB(conf.GetConf())
	if err != nil {
		panic(err)
	}

	if err = db.InitDb(); err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: handlers.NewRouter(db),
	}

	log.Print("Starting server http service...")
	log.Printf("Listening on port %v", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
