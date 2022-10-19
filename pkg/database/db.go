package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dylan-dinh/api-gin/pkg/conf"
	"github.com/dylan-dinh/api-gin/pkg/model"
	"github.com/jmoiron/modl"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
)

type registerer interface {
	Register(*modl.DbMap)
}

// DB describe the database
type DB struct {
	conf   *conf.Conf
	logger *log.Logger
	*modl.DbMap
}

func NewDB(conf *conf.Conf) (*DB, error) {
	database, err := sql.Open(conf.DriverName, conf.ConnectionStr)
	if err != nil {
		return nil, err
	}

	return &DB{
		conf:   conf,
		DbMap:  modl.NewDbMap(database, modl.PostgresDialect{}),
		logger: log.New(os.Stdout, "[DB]: ", 2),
	}, nil
}

func (db *DB) InitDb() error {
	models := []registerer{
		&model.Site{},
		&model.Engine{},
	}

	for _, m := range models {
		m.Register(db.DbMap)
	}

	db.logger.Print("Creating table...")
	if err := db.DbMap.CreateTablesIfNotExists(); err != nil {
		return err
	}

	if err := db.addSite(db.DbMap); err != nil {
		return err
	}

	if err := addTrigger(db.DbMap); err != nil {
		return err
	}
	return nil
}

func addTrigger(dbMap *modl.DbMap) error {
	query := `CREATE OR REPLACE FUNCTION process_site_max_power() RETURNS trigger AS $site_max_power$
	DECLARE
	total INTEGER;
	maxPower INTEGER;
	BEGIN
	SELECT sum(rated_capacity) FROM engines WHERE site_id = NEW.site_id INTO total;
	SELECT max_power FROM sites WHERE id = NEW.site_id INTO maxPower;
	
	IF total > maxPower THEN
	RAISE EXCEPTION 'rated capacity of whole engine power cannot exceed %', maxPower;
	END IF;
	RETURN NEW;
	END;
	$site_max_power$ LANGUAGE plpgsql;
	
	CREATE OR REPLACE TRIGGER site_max_power
	AFTER INSERT OR UPDATE ON engines FOR EACH ROW EXECUTE FUNCTION process_site_max_power();
	`

	_, err := dbMap.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) addSite(dbmap *modl.DbMap) error {
	var site model.Site
	var siteMaxPower int64

	err := dbmap.SelectOne(&site, `SELECT * FROM sites WHERE name = $1`, db.conf.SiteName)
	// err != nil means nothing is found, so we insert it
	if err != nil {
		if siteMaxPower, err = strconv.ParseInt(db.conf.SiteMaxPower, 10, 64); err != nil {
			return errors.New(fmt.Sprintf("error while parsing site max power : %s", err.Error()))
		}
		site.Name = db.conf.SiteName
		site.MaxPower = siteMaxPower

		err = dbmap.Insert(&site)
		if err != nil {
			return errors.New(fmt.Sprintf("error while creating site : %s", err.Error()))
		}

		// found so we just update the max power in case
	} else {
		if siteMaxPower, err = strconv.ParseInt(db.conf.SiteMaxPower, 10, 64); err != nil {
			return errors.New(fmt.Sprintf("error while parsing site max power : %s", err.Error()))
		}
		site.MaxPower = siteMaxPower
		if _, err = dbmap.Update(&site); err != nil {
			return errors.New(fmt.Sprintf("error while updating site : %s", err.Error()))
		}
	}
	return nil
}

func healthCheck(db *DB) error {
	err := db.DbMap.Db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Start() error {
	var err error

	db.logger.Print("Starting database service...")
	err = db.InitDb()
	if err != nil {
		return err
	}

	db.logger.Print("Pinging database...")
	err = healthCheck(db)
	if err != nil {
		return err
	}
	db.logger.Print("Connected.")

	return nil
}
