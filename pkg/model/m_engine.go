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
	SiteId        int64  `json:"site_id" db:"site_id" `
	Name          string `json:"name" db:"name" `
	Type          string `json:"type" db:"type"`
	RatedCapacity int64  `json:"rated_capacity" db:"rated_capacity" `
}

func (s Engine) Register(db *modl.DbMap) {
	engine := db.AddTableWithName(s, "engines").SetKeys(true, "id")
	engine.ColMap("site_id").SetSqlType("smallserial")
	engine.ColMap("name").SetSqlType("varchar(50) NOT NULL DEFAULT ''")
	engine.ColMap("type").SetSqlType("varchar(50) NOT NULL DEFAULT ''")
	engine.ColMap("rated_capacity").SetSqlType("serial")
}

func validateData(s *Engine, v *Validator) {
	if s.Type != TypeFurnace && s.Type != TypeCompressor &&
		s.Type != TypeChiller && s.Type != TypeRollingMill {
		v.AddError("type", fmt.Sprintf("%s is not a valid type of engine", s.Type))
	}

	// to do valid site id exist

	if s.Name == "" {
		v.AddError("name", "shouldn't be empty")
	}

	if s.RatedCapacity <= 0 {
		v.AddError("rated_capacity", "shouldn't be equal or inferior to 0")
	}

}

func (s *Engine) Validate(ex modl.SqlExecutor) ValidationErrors {
	v := NewValidator(ex)

	validateData(s, v)

	return v.Errors()
}

func (s *Engine) PreInsert(ex modl.SqlExecutor) error {
	if err := s.Validate(ex); len(err) != 0 {
		return err
	}

	return nil
}

func (s *Engine) PreUpdate(ex modl.SqlExecutor) error {
	if err := s.Validate(ex); len(err) != 0 {
		return err
	}

	return nil
}
