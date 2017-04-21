package db_user_r

import (
	"github.com/go-xorm/xorm"

	"github.com/mamoroom/echo-mvc/src/models/db"
	"github.com/mamoroom/echo-mvc/src/models/db/db_conns"
)

func New() db_conns.DbConn {
	return DbConnUserR{
		db: db.DbUserR,
	}
}

type DbConnUserR struct {
	db *xorm.Engine
	db_conns.DbConnBase
}

func (db_conn_user_r DbConnUserR) Engine() *xorm.Engine {
	return db_conn_user_r.db
}

func (db_conn_user_r DbConnUserR) GetDbConnName() string {
	return "user_r"
}
