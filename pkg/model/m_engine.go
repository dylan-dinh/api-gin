package model

import (
	"fmt"
	"github.com/jmoiron/modl"
)

const (
	TypeFurnace     = "furnace"
	TypeCompressor  = "compressor"
	TypeChiller     = "chiller"
	TypeRollingMill = "rolling mill"
)

type Engine struct {
	// PK, auto increment
	Id            int64  `json:"id" db:"id" binding:"-"`
	SiteId        int64  `json:"site_id" db:"site_id" binding:"required"`
	Name          string `json:"name" db:"name" binding:"required"`
	Type          string `json:"type" db:"type" binding:"required"`
	RatedCapacity int64  `json:"rated_capacity" db:"rated_capacity" binding:"required"`
}

func (e Engine) Register(db *modl.DbMap) {
	engine := db.AddTableWithName(e, "engines").SetKeys(true, "id")
	engine.ColMap("site_id").SetSqlType("smallserial")
	engine.ColMap("name").SetSqlType("varchar(50) NOT NULL DEFAULT ''")
	engine.ColMap("type").SetSqlType("varchar(50) NOT NULL DEFAULT ''")
	engine.ColMap("rated_capacity").SetSqlType("serial")
}

func validateEngineData(e *Engine, v *Validator) {
	if e.Type != TypeFurnace && e.Type != TypeCompressor &&
		e.Type != TypeChiller && e.Type != TypeRollingMill {
		v.AddError("type", fmt.Sprintf("%s is not a valid type of engine", e.Type))
	}

	// to do valid site id exist
	v.Exists("sites", "id", e.SiteId, "site_id", fmt.Sprintf("site with id '%d' does not exist", e.SiteId))

	if e.Name == "" {
		v.AddError("name", "shouldn't be empty")
	}

	if e.RatedCapacity <= 0 {
		v.AddError("rated_capacity", "shouldn't be equal or inferior to 0")
	}
}

func (e *Engine) Validate(ex modl.SqlExecutor) ValidationErrors {
	v := NewValidator(ex)

	validateEngineData(e, v)

	return v.Errors()
}

func (e *Engine) PreInsert(ex modl.SqlExecutor) error {
	if err := e.Validate(ex); len(err) != 0 {
		return err
	}

	return nil
}

func (e *Engine) PreUpdate(ex modl.SqlExecutor) error {
	if err := e.Validate(ex); len(err) != 0 {
		return err
	}

	return nil
}
