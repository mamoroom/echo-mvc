package db_conns

import (
	"github.com/go-xorm/xorm"
)

type DbConn interface {
	Engine() *xorm.Engine
	NewSession() *xorm.Session
	GetDbConnName() string
	ClearCache(beans ...interface{}) error
}

type DbConnBase struct{}

func (db_c_base DbConnBase) GetDbConnName() string {
	panic("Must implement as child class")
}

func (db_c_base DbConnBase) Engine() *xorm.Engine {
	panic("Must implement as child class")
}

func (db_c_base DbConnBase) NewSession() *xorm.Session {
	panic("Must implement as child class")
}

func (db_c_base DbConnBase) ClearCache(beans ...interface{}) error {
	panic("Must implement as child class")
}
