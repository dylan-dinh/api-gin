package model

import "github.com/jmoiron/modl"

type Site struct {
	// PK, auto increment
	Id   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (s Site) Register(db *modl.DbMap) {
	site := db.AddTableWithName(s, "sites").SetKeys(true, "id")
	site.ColMap("name").SetSqlType("varchar(50) NOT NULL DEFAULT ''")

}
