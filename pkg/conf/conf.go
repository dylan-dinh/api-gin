package conf

import (
	"flag"
	"fmt"
	"github.com/graniticio/inifile"
	"log"
	"os"
)

type Conf struct {
	SiteMaxPower  string
	DriverName    string `default:"postgres"`
	ConnectionStr string `description:"Connection string to postgres"`
}

// HandleArgs check if conf file is passed
func HandleArgs() {
	if len(os.Args) < 2 {
		log.Fatal("program needs at least configuration file to work properly")
	}
}

// GetConf is it better to use chan and log every error instead of stopping to the first ?
// read me will explain how to configure ini file so ?
func GetConf() *Conf {
	logger := log.New(os.Stdout, "[CONFIG]:  ", 2)

	file := flag.String("c", "", "configuration file")
	flag.Parse()

	ic, err := inifile.NewIniConfigFromPath(*file)
	if err != nil {
		logger.Fatal("error loading configuration file")
	}

	if !ic.SectionExists("database") || !ic.SectionExists("site") {
		logger.Fatal("section 'database' and/or section 'site' are missing")
	}

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

	maxPower, err := ic.Value("site", "MaxPower")
	if err != nil {
		logger.Fatal("property MaxPower is missing in database section")
	}

	return &Conf{
		SiteMaxPower: maxPower,
		DriverName:   "postgres",
		ConnectionStr: fmt.Sprintf("host=%s port=%d dbname=%s user=%s password='%s' sslmode=%s",
			host, port, name, username, pswd, "disable"),
	}
}
