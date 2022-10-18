package main

import (
	"github.com/dylan-dinh/api-gin/pkg/conf"
	"github.com/dylan-dinh/api-gin/pkg/database"
	"github.com/dylan-dinh/api-gin/pkg/handlers"
	"log"
)

type Services struct {
	Db *database.DB
	Ap *handlers.AppServer
}

func (s *Services) startServices() {
	err := s.Db.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

	s.Ap.Db = s.Db.GetDB()

	err = s.Ap.Start()
	if err != nil {
		log.Fatal(err.Error())
	}

}

func main() {
	var db database.DB
	var ap handlers.AppServer

	conf.HandleArgs()

	db.Conf = conf.GetConf()

	s := Services{Db: &db, Ap: &ap}

	s.startServices()
}
