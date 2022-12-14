package conf

import (
	"flag"
	"fmt"
	"github.com/graniticio/inifile"
	"log"
	"os"
)

// Conf holds configuration for db and site
type Conf struct {
	SiteMaxPower  string
	SiteName      string
	DriverName    string `default:"postgres"`
	ConnectionStr string `description:"Connection string to postgres"`
}

// GetConf retrieve conf from conf file
func GetConf() *Conf {
	if len(os.Args) < 2 {
		log.Fatal("program needs at least configuration file to work properly")
	}

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

	return &Conf{
		SiteMaxPower: maxPower,
		SiteName:     siteName,
		DriverName:   "postgres",
		ConnectionStr: fmt.Sprintf("host=%s port=%d dbname=%s user=%s password='%s' sslmode=%s",
			host, port, name, username, pswd, "disable"),
	}
}
