package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/xorm/migrate"

	"github.com/mamoroom/echo-mvc/src/models/db"
	"github.com/mamoroom/echo-mvc/src/models/entity"

	_ "fmt"
	"log"
)

var (
	migrations = []*migrate.Migration{
		{
			ID: "20170419-6",
			Migrate: func(tx *xorm.Engine) error {
				return tx.Sync2(&entity.User{})
			},
			Rollback: func(tx *xorm.Engine) error {
				return tx.DropTables(&entity.User{})
			},
		},
	}
)

func main() {
	defer db.DbUserW.Close()

	if err := db.DbUserW.DB().Ping(); err != nil {
		log.Fatal(err)
	}
	m := migrate.New(db.DbUserW, migrate.DefaultOptions, migrations)
	err := m.Migrate()
	if err != nil {
		log.Fatal(err)
	}
}
