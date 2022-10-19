package test

import (
	"flag"
	"fmt"
	"github.com/dylan-dinh/api-gin/pkg/conf"
	"github.com/dylan-dinh/api-gin/pkg/database"
	"github.com/graniticio/inifile"
	"log"
	"os"
)

var confFile = flag.String("conf", "", "")

func ClearDB(db *database.DB) {
	db.Db.Exec("truncate sites")
	db.Db.Exec("truncate engines")
}

func SetUpDbTest() (*database.DB, *conf.Conf) {
	cf := getConf()

	db, err := database.NewDB(cf)
	if err != nil {
		log.Fatal(err)
	}

	err = db.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	return db, cf
}

func getConf() *conf.Conf {
	logger := log.New(os.Stdout, "[CONFIG]:  ", 2)

	ic, err := inifile.NewIniConfigFromPath(*confFile)
	if err != nil {
		logger.Fatal(err.Error())
	}

	if !ic.SectionExists("database") || !ic.SectionExists("site") {
		logger.Fatal("section 'database' and/or section 'site' are missing")
	}

	// db conf
	host, err := ic.Value("database", "Host")
	if err != nil {
		logger.Fatal("property Host is missing in database section")
	}

	name, err := ic.Value("database", "Name")
	if err != nil {
		logger.Fatal("property Name is missing in database section")
	}

	port, err := ic.ValueAsInt64("database", "Port")
	if err != nil {
		logger.Fatal("property Port is missing in database section")
	}

	username, err := ic.Value("database", "User")
	if err != nil {
		logger.Fatal("property User is missing in database section")
	}

	pswd, err := ic.Value("database", "Password")
	if err != nil {
		logger.Fatal("property Password is missing in database section")
	}

	// site conf
	maxPower, err := ic.Value("site", "MaxPower")
	if err != nil {
		logger.Fatal("property MaxPower is missing in site section")
	}

	siteName, err := ic.Value("site", "SiteName")
	if err != nil {
		logger.Fatal("property SiteName is missing in site section")
	}

	return &conf.Conf{
		SiteMaxPower: maxPower,
		SiteName:     siteName,
		DriverName:   "postgres",
		ConnectionStr: fmt.Sprintf("host=%s port=%d dbname=%s user=%s password='%s' sslmode=%s",
			host, port, name, username, pswd, "disable"),
	}
}
