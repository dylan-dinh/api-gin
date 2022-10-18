package handlers

import (
	"github.com/dylan-dinh/api-gin/pkg/handlers/engines"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/modl"
	"log"
	"net/http"
	"os"
)

type AppServer struct {
	Logger *log.Logger
	Db     *modl.DbMap
}

func (ap *AppServer) Start() error {
	router := gin.Default()
	ap.Logger = log.New(os.Stdout, "[SERVER]:  ", 2)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	engines.NewHandler(&engines.Config{
		Ngine: router,
		Db:    ap.Db,
	})

	ap.Logger.Print("Starting server http service...")
	ap.Logger.Printf("Listening on port %v", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
