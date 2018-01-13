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
	AutoFill(userName string)
	GetNumber() int
	Kind() string
	GetName() string
	ShortDesc() string
	Desc() string
	GetSequence() string
	OwnerFieldName() string
	GetCoreLinks() *CoreLinks
	GetCoreInfoSections() []InfoSection
	GetSequenceInfo() *SequenceInfo
	GetSupplementalFields() []FieldDef
	AsResourceDef() string
}

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *Model) model() Model                       { return *m }
func (m *Model) GetID() uint                        { return m.ID }
func (m *Model) AutoFill(userName string)           {}
func (m *Model) GetNumber() int                     { return int(m.ID) }
func (m *Model) Kind() string                       { return "" }
func (m *Model) GetName() string                    { return "" }
func (m *Model) ShortDesc() string                  { return "" }
func (m *Model) Desc() string                       { return "" }
func (m *Model) GetSequence() string                { return "" }
func (m *Model) OwnerFieldName() string             { return "name" }
func (m *Model) GetCoreLinks() *CoreLinks           { return nil }
func (m *Model) GetCoreInfoSections() []InfoSection { return nil }
func (m *Model) GetSequenceInfo() *SequenceInfo     { return nil }
func (m *Model) GetSupplementalFields() []FieldDef  { return nil }
func (m *Model) AsResourceDef() string              { return "" }

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
	case "rnaiclone", "rnaiclones":
		return &RNAiClone{}
	case "seqlib", "seqlibs":
		return &SeqLib{}
	default:
		return &Model{}
	}
}

type EntityQueryIterator struct {
	query        *gorm.DB
	buffer       []Entity
	bufferOffset int
	pageOffset   int
	cls          string
}

func (eqi *EntityQueryIterator) HasNext() bool {
	if eqi.buffer != nil && eqi.bufferOffset < len(eqi.buffer) {
		return true
	}
	test := Empty(eqi.cls)
	eqi.query.Offset(eqi.pageOffset + eqi.bufferOffset).Limit(1).First(test)
	return test.GetID() > 0
}

func (eqi *EntityQueryIterator) Next() Entity {
	pageSize := 100
	if eqi.buffer == nil || eqi.bufferOffset >= len(eqi.buffer) {
		eqi.pageOffset += eqi.bufferOffset
		eqi.bufferOffset = 0
		page := RunQuery(eqi.cls, eqi.query.Offset(eqi.pageOffset).Limit(pageSize))
		eqi.buffer = page
	}
	result := eqi.buffer[eqi.bufferOffset]
	eqi.bufferOffset += 1
	return result
}

func RunQueryLazy(cls string, db *gorm.DB) *EntityQueryIterator {
	return &EntityQueryIterator{
		query:        db,
		buffer:       nil,
		bufferOffset: 0,
		pageOffset:   0,
		cls:          cls,
	}
}

func RunQuery(cls string, db *gorm.DB) []Entity {
	entityResult := []Entity{}
	switch cls {
	case "plasmid", "plasmids":
		res := []Plasmid{}
		db.Find(&res)
		for i, _ := range res {
			entityResult = append(entityResult, &res[i])
		}
	case "oligo", "oligos":
		res := []Oligo{}
		db.Find(&res)
		for i, _ := range res {
			entityResult = append(entityResult, &res[i])
		}
	case "line", "lines":
		res := []Line{}
		db.Find(&res)
		for i, _ := range res {
			entityResult = append(entityResult, &res[i])
		}
	case "sample", "samples":
		res := []Sample{}
		db.Find(&res)
		for i, _ := range res {
			entityResult = append(entityResult, &res[i])
		}
	case "bacterium", "bacteria":
		res := []Bacterium{}
		db.Find(&res)
		for i, _ := range res {
			entityResult = append(entityResult, &res[i])
		}
	case "yeaststrain", "yeaststrains":
		res := []Yeaststrain{}
		db.Find(&res)
		for i, _ := range res {
			entityResult = append(entityResult, &res[i])
		}
	case "user", "users":
		res := []User{}
		db.Find(&res)
		for i, _ := range res {
			entityResult = append(entityResult, &res[i])
		}
	case "antibody", "antibodies":
		res := []Antibody{}
		db.Find(&res)
		for i, _ := range res {
			entityResult = append(entityResult, &res[i])
		}
	case "rnaiclone", "rnaiclones":
		res := []RNAiClone{}
		db.Find(&res)
		for i, _ := range res {
			entityResult = append(entityResult, &res[i])
		}
	case "seqlib", "seqlibs":
		res := []SeqLib{}
		db.Find(&res)
		for i, _ := range res {
			entityResult = append(entityResult, &res[i])
		}
	default:
		res := []Model{}
		db.Find(&res)
		for i, _ := range res {
			entityResult = append(entityResult, &res[i])
		}
	}
	return entityResult
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
	db.Where("id > ?", id).First(e)
	return e
}

func Prev(e Entity, id string) Entity {
	db.Where("id < ?", id).Last(e)
	return e
}

func NextAvailableNumber(e Entity) int {
	// TODO(colin): possible race here with multiple people creating items of the
	// same type at the same time.
	db.Last(e)
	return e.GetNumber() + 1
}

func GetByID(e Entity, id int) {
	db.First(e, id)
}

func Create(e Entity) {
	db.Create(e)
}

func init() {
	dbURL := env.DbURL
	if env.Dev {
		dbURL = "dbname=labdb sslmode=disable"
	}
	pg, err := gorm.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	db = pg

	db.AutoMigrate(&SeqLib{})
	db.AutoMigrate(&RNAiClone{})
	db.LogMode(env.DebugDB)
}

func Shutdown() {
	if db != nil {
		db.Close()
	}
}

func Db() *gorm.DB {
	return db
}

func IsImplemented(modelType string) bool {
	return false
	//  return (modelType == "rnai_clone") || (modelType == "seq_lib")
}
