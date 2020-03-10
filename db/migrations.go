package db

import (
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/xorm/migrate"
)

var migrations = []*migrate.Migration{
	{
		ID: "initial_202003101829",
		Migrate: func(tx *xorm.Engine) error {
			return tx.Sync2(&DbOffer{})
		},
		Rollback: func(tx *xorm.Engine) error {
			return tx.DropTables(&DbOffer{})
		},
	},
}
