package models

import (
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Enables the postgres driver for gorm.

	"labdb.org/labdb/env"
)

var db *gorm.DB

type Entity interface {
	model() Model
	GetID() uint
}

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Model) model() Model {
	return *m
}

func (m *Model) GetID() uint {
	return m.ID
}

type Plasmid struct {
	Model
}

type Oligo struct {
	Model
}

type Line struct {
	Model
}

type Sample struct {
	Model
}

type Bacterium struct {
	Model
}

type Yeaststrain struct {
	Model
}

type User struct {
	Model
}

type Antibody struct {
	Model
}

func Empty(cls string) Entity {
	switch cls {
	case "plasmid", "plasmids":
		return &Plasmid{}
	case "oligo", "oligos":
		return &Oligo{}
	case "line", "lines":
		return &Line{}
	case "sample", "samples":
		return &Sample{}
	case "bacterium", "bacteria":
		return &Bacterium{}
	case "yeaststrain", "yeaststrains":
		return &Yeaststrain{}
	case "user", "users":
		return &User{}
	case "antibody", "antibodies":
		return &Antibody{}
	default:
		return &Model{}
	}
}

func NextID(cls string, currID string) string {
	ent := Empty(cls)
	id := Next(ent, currID).GetID()
	if id != 0 {
		return strconv.FormatUint(uint64(id), 10)
	}
	return currID
}

func PrevID(cls string, currID string) string {
	ent := Empty(cls)
	id := Prev(ent, currID).GetID()
	if id != 0 {
		return strconv.FormatUint(uint64(id), 10)
	}
	return currID
}

func Next(e Entity, id string) Entity {
	db.Where("id > ?", id).Order("id asc").First(e)
	return e
}

func Prev(e Entity, id string) Entity {
	db.Where("id < ?", id).Order("id desc").First(e)
	return e
}

func Init() {
	dbURL := env.DbURL
	if env.Dev {
		dbURL = "dbname=labdb sslmode=disable"
	}
	pg, err := gorm.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	db = pg
}

func Shutdown() {
	if db != nil {
		db.Close()
	}
}

func Db() *gorm.DB {
	return db
}
