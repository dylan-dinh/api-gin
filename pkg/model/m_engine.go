package model

import "github.com/jmoiron/modl"

type Engine struct {
	// PK, auto increment
	Id            int64  `json:"id" db:"id"`
	SiteId        int64  `json:"site_id" db:"site_id"`
	Name          string `json:"name" db:"name"`
	Type          string `json:"type" db:"type"`
	RatedCapacity int64  `json:"rated_capacity" db:"rated_capacity"`
}

func (s Engine) Register(db *modl.DbMap) {
	engine := db.AddTableWithName(s, "engines").SetKeys(true, "id")
	engine.ColMap("site_id").SetSqlType("smallserial")
	engine.ColMap("name").SetSqlType("varchar(50) NOT NULL DEFAULT ''")
	engine.ColMap("type").SetSqlType("varchar(50) NOT NULL DEFAULT ''")
	engine.ColMap("rated_capacity").SetSqlType("serial")

}
