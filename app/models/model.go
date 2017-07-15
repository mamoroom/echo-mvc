package models

import (
	"github.com/go-xorm/core"

	"github.com/mamoroom/echo-mvc/app/config"
	"github.com/mamoroom/echo-mvc/app/models/dbh"

	_ "fmt"
)

var conf = config.Conf

func NewBase(dbh *dbh.DbHandler) *BaseModel {
	return &BaseModel{
		Dbh: dbh,
	}
}

type BaseModel struct {
	Dbh *dbh.DbHandler
}

// common //
func (b *BaseModel) GetTableName(bean interface{}) string {
	return b.Dbh.GetTableName(bean)
}

func (b *BaseModel) ClearCache(beans ...interface{}) error {
	return b.Dbh.ClearCache(beans...)
}

// query //

// GET
func (b *BaseModel) GetById(id interface{}, bean interface{}) (bool, error) {
	return b.Dbh.Handle().Id(id).Get(bean)
}

func (b *BaseModel) GetByPk(pk *core.PK, bean interface{}) (bool, error) {
	return b.Dbh.Handle().Id(pk).Get(bean)
}

func (b *BaseModel) GetBy(statement string, value interface{}, beans interface{}) (bool, error) {
	return b.Dbh.Handle().Where(statement, value).Get(beans)
}

// Find
func (b *BaseModel) Find(beans interface{}) error {
	return b.Dbh.Handle().Find(beans)
}

func (b *BaseModel) FindDescBy(col string, beans interface{}) error {
	return b.Dbh.Handle().Desc(col).Find(beans)
}

func (b *BaseModel) FindAscBy(col string, beans interface{}) error {
	return b.Dbh.Handle().Asc(col).Find(beans)
}

func (b *BaseModel) FindBy(statement string, value interface{}, beans interface{}) error {
	return b.Dbh.Handle().Where(statement, value).Find(beans)
}

func (b *BaseModel) FindByDescOf(statement string, value interface{}, col string, beans interface{}) error {
	return b.Dbh.Handle().Where(statement, value).Desc(col).Find(beans)
}

func (b *BaseModel) FindByAscOf(statement string, value interface{}, col string, beans interface{}) error {
	return b.Dbh.Handle().Where(statement, value).Asc(col).Find(beans)
}

// UPDATE
func (b *BaseModel) UpdateById(id interface{}, bean interface{}, set_cols ...string) (int64, error) {
	return b.Dbh.Handle().Id(id).Cols(set_cols...).Update(bean)
}

func (b *BaseModel) UpdateByPk(pk *core.PK, bean interface{}, set_cols ...string) (int64, error) {
	return b.Dbh.Handle().Id(pk).Cols(set_cols...).Update(bean)
}

// INSERT
func (b *BaseModel) Insert(bean interface{}, omit_cols ...string) (int64, error) {
	return b.Dbh.Handle().Omit(omit_cols...).Insert(bean)
}

// query opt //
func (b *BaseModel) SetNullableCols(nullable_cols ...string) {
	b.Dbh.SetNullable(nullable_cols...)
}

// EXEC raw statement //
func (b *BaseModel) Exec(sqlStr string, args ...interface{}) (int64, error) {
	result, err := b.Dbh.Exec(sqlStr, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
