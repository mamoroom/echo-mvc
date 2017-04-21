package models

import (
	"github.com/mamoroom/echo-mvc/src/models/db/db_conns"
	"github.com/mamoroom/echo-mvc/src/models/dbh"
)

func NewBase(db_conn db_conns.DbConn) *BaseModel {
	return &BaseModel{
		Dbh: dbh.New(db_conn),
	}
}

type BaseModel struct {
	Dbh *dbh.DbHandler
}
