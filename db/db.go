package db

import (
	"ets2-sync/dlc"
	"ets2-sync/global"
	"fmt"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/xorm/migrate"
	_ "github.com/lib/pq"
	"os"
	"time"
	"xorm.io/core"
)

var db *xorm.Engine

type DbOffer struct {
	Id                 string `xorm:"pk text"`
	RequiredDlc        dlc.Dlc
	SourceCompany      string `xorm:"text"`
	Target             string
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
}

func InitializeDb() error {
	if db != nil {
		return nil
	}

	getConnectionString := func(dbName string) string {
		return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
			os.Getenv("dbHost"), os.Getenv("dbPort"),
			os.Getenv("dbUser"), os.Getenv("dbPassword"), dbName)
	}

	engine, err := xorm.NewEngine("postgres", getConnectionString("postgres"))

	if err != nil {
		return err
	}

	realDbName := os.Getenv("dbName")

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
