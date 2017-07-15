package db_user_w

import (
	"github.com/go-xorm/xorm"

	"github.com/mamoroom/echo-mvc/app/models/db"
	"github.com/mamoroom/echo-mvc/app/models/db/db_conns"
)

func New() db_conns.DbConn {
	return DbConnUserW{
		db: db.DbUserW,
	}
}

type DbConnUserW struct {
	db *xorm.Engine
	db_conns.DbConnBase
}

func (db_conn_user_w DbConnUserW) Engine() *xorm.Engine {
	return db_conn_user_w.db
}

func (db_conn_user_w DbConnUserW) NewSession() *xorm.Session {
	return db_conn_user_w.db.NewSession()
}

func (db_conn_user_w DbConnUserW) GetDbConnName() string {
	return "user_w"
}
