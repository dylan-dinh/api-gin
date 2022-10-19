package model

import "github.com/jmoiron/modl"

type Site struct {
	// PK, auto increment
	Id       int64  `json:"id" db:"id" binding:"-"`
	Name     string `json:"name" db:"name" binding:"required"`
	MaxPower int64  `json:"max_power" db:"max_power" binding:"required"`
}

func (s Site) Register(db *modl.DbMap) {
	site := db.AddTableWithName(s, "sites").SetKeys(true, "id")
	site.ColMap("name").SetSqlType("varchar(50) NOT NULL DEFAULT ''")
	site.ColMap("max_power").SetSqlType("bigserial")
}

func (s *Site) Validate(ex modl.SqlExecutor) ValidationErrors {
	v := NewValidator(ex)

	if s.MaxPower <= 0 {
		v.AddError("max_power", "shouldn't be equal or inferior to 0")
	}

	return v.Errors()
}

func (s *Site) PreInsert(ex modl.SqlExecutor) error {
	if err := s.Validate(ex); len(err) != 0 {
		return err
	}

	return nil
}

func (s *Site) PreUpdate(ex modl.SqlExecutor) error {
	if err := s.Validate(ex); len(err) != 0 {
		return err
	}

	return nil
}
