package db

import (
	"ets2-sync/dlc"
	"ets2-sync/global"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/xorm/migrate"
	_ "github.com/lib/pq"
	"github.com/mitchellh/hashstructure"
	"os"
	"strconv"
	"time"
	"xorm.io/core"
)

var db *xorm.Engine

type DbOffer struct {
	Id                 string `xorm:"pk text"`
	RequiredDlc        dlc.Dlc
	SourceCompany      string `xorm:"text"`
	Target             string
	ExpirationTime     string
	Urgency            string
	ShortestDistanceKm string
	FerryTime          string
	FerryPrice         string
	Cargo              string
	CompanyTruck       string
	TrailerVariant     string
	TrailerDefinition  string
	UnitsCount         string
	FillRatio          string
	TrailerPlace       string
	NameParam          string // nameParam
}

func (o *DbOffer) CalculateHash() string {
	hash, err := hashstructure.Hash(struct {
		S string
		T string
		C string
	}{o.SourceCompany, o.Target, o.Cargo}, nil)

	if err != nil {
		return ""
	}

	return strconv.FormatUint(hash, 10)
}

func InitializeDb() error {
	if db != nil {
		return nil
	}

	getConnectionString := func(dbName string) string {
		return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
			os.Getenv("db-host"), os.Getenv("db-port"),
			os.Getenv("db-user"), os.Getenv("db-password"), dbName)
	}

	engine, err := xorm.NewEngine("postgres", getConnectionString("postgres"))

	if err != nil {
		return err
	}

	realDbName := os.Getenv("db-name")

	if engine != nil {
		_, _ = engine.Exec(fmt.Sprintf("CREATE DATABASE %s;", realDbName))

		_ = engine.Close()
	}

	engine, err = xorm.NewEngine("postgres", getConnectionString(realDbName))

	if err != nil {
		return err
	}

	db = engine

	engine.SetSchema("public")
	engine.SetMapper(core.SnakeMapper{})
	engine.SetTZLocation(time.UTC)
	engine.SetTZDatabase(time.UTC)

	if global.IsDebug {
		engine.ShowSQL(true)
	}

	m := migrate.New(engine, &migrate.Options{
		TableName:    "migrations",
		IDColumnName: "id",
	}, migrations)

	err = m.Migrate()

	if err != nil {
		return err
	}

	return err
}

func GetDb() *xorm.Engine {
	return db
}
