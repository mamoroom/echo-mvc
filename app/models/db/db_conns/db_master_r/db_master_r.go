package db_master_r

import (
	"github.com/go-xorm/xorm"

	"github.com/mamoroom/echo-mvc/app/models/db"
	"github.com/mamoroom/echo-mvc/app/models/db/db_conns"
)

func New() db_conns.DbConn {
	return DbConnMasterR{
		db: db.DbMasterR,
	}
}

type DbConnMasterR struct {
	db *xorm.Engine
	db_conns.DbConnBase
}

func (db_conn_master_r DbConnMasterR) Engine() *xorm.Engine {
	return db_conn_master_r.db
}

func (db_conn_master_r DbConnMasterR) GetDbConnName() string {
	return "master_r"
}

func (db_conn_master_r DbConnMasterR) ClearCache(beans ...interface{}) error {
	return db_conn_master_r.db.ClearCache(beans...)
}
