package main

import (
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/xorm/migrate"
)

var migrations = []*migrate.Migration{
	{
		ID: "initial_202003101829",
		Migrate: func(tx *xorm.Engine) error {
			return tx.Sync2(&dbOffer{})
		},
		Rollback: func(tx *xorm.Engine) error {
			return tx.DropTables(&dbOffer{})
		},
	},
}
