package main

import (
	"github.com/dylan-dinh/api-gin/pkg/conf"
	"github.com/dylan-dinh/api-gin/pkg/database"
	"github.com/dylan-dinh/api-gin/pkg/handlers"
	"log"
	"os"
)

// handleArgs check if conf file is passed
func handleArgs() {
	if len(os.Args) < 2 {
		log.Fatal("program needs at least configuration file to work properly")
	}
}

type Services struct {
	Db *database.DB
	Ap *handlers.AppServer
}

func (s *Services) startServices() {

	s.Db.Start()
	s.Ap.Db = s.Db.GetDB()
	s.Ap.Start()

}

func main() {
	var db database.DB
	var ap handlers.AppServer

	handleArgs()

	db.Conf = conf.GetConf()

	s := Services{Db: &db, Ap: &ap}

	s.startServices()
}
