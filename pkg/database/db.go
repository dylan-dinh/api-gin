package database

import (
	"database/sql"
	"github.com/dylan-dinh/api-gin/pkg/conf"
	"github.com/dylan-dinh/api-gin/pkg/model"
	"github.com/jmoiron/modl"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type registerer interface {
	Register(*modl.DbMap)
}

// DB describe the database
type DB struct {
	Conf   *conf.Conf
	DbMap  *modl.DbMap
	Logger *log.Logger
}

// Close cl
func (db *DB) Close() {
	err := db.DbMap.Db.Close()
	if err != nil {
		db.Logger.Printf("error closing db : %s", err.Error())
	} else {
		db.Logger.Print("db closed")
	}
}

func initDb(db *DB) (*modl.DbMap, error) {
	database, err := sql.Open(db.Conf.DriverName, db.Conf.ConnectionStr)
	if err != nil {
		return nil, err
	}

	dbMap := modl.NewDbMap(database, modl.PostgresDialect{})

	models := []registerer{
		&model.Site{},
		&model.Engine{},
	}

	for _, m := range models {
		m.Register(dbMap)
	}

	db.Logger.Print("Creating table...")
	if err = dbMap.CreateTablesIfNotExists(); err != nil {
		return nil, err
	}

	return dbMap, nil
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
	db.Logger = log.New(os.Stdout, "[DB]: ", 2)

	db.Logger.Print("Starting database service...")
	db.DbMap, err = initDb(db)
	if err != nil {
		return err
	}

	db.Logger.Print("Pinging database...")
	err = healthCheck(db)
	if err != nil {
		return err
	}
	db.Logger.Print("Connected.")

	return nil
}

func (db *DB) GetDB() *modl.DbMap {
	return db.DbMap
}
