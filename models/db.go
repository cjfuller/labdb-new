package models

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Enables the postgres driver for gorm.

	"labdb.org/labdb/env"
)

var dbByLab = map[string]*gorm.DB{}

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Plasmid struct {
	Model
	Name     string
	Alias    string
	Sequence string
}

func Init() {
	dbURL := "dbname=labdb"
	if env.Dev {
		dbURL += " sslmode=disable"
	}
	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	dbByLab["dev"] = db
}

func Shutdown() {
	for _, db := range dbByLab {
		db.Close()
	}
}

func Db(lab string) *gorm.DB {
	return dbByLab[lab]
}
