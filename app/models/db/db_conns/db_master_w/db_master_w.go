package db_master_w

import (
	"github.com/go-xorm/xorm"

	"github.com/mamoroom/echo-mvc/app/models/db"
	"github.com/mamoroom/echo-mvc/app/models/db/db_conns"
)

func New() db_conns.DbConn {
	return DbConnMasterW{
		db: db.DbMasterW,
	}
}

type DbConnMasterW struct {
	db *xorm.Engine
	db_conns.DbConnBase
}

func (db_conn_master_w DbConnMasterW) Engine() *xorm.Engine {
	return db_conn_master_w.db
}

func (db_conn_master_w DbConnMasterW) NewSession() *xorm.Session {
	return db_conn_master_w.db.NewSession()
}

func (db_conn_master_w DbConnMasterW) GetDbConnName() string {
	return "master_w"
}
