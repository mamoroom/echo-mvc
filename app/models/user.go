package models

import (
	"github.com/mamoroom/echo-mvc/app/models/db/db_conns/db_user_r"
	"github.com/mamoroom/echo-mvc/app/models/dbh"
	"github.com/mamoroom/echo-mvc/app/models/entity"
)

func NewUserR() *UserModel {
	dbh := dbh.New(db_user_r.New())
	return NewUser(dbh)
}

func NewUser(dbh *dbh.DbHandler) *UserModel {
	b := NewBase(dbh)
	return &UserModel{
		BaseModel:          *b,
		TableName:          b.GetTableName(entity.User{}),
		update_target_cols: map[string]string{},
	}
}

type UserModel struct {
	TableName string
	BaseModel
	entity             *entity.User
	update_target_cols map[string]string
}

func (m *UserModel) SetEntity(bean *entity.User) {
	m.entity = bean
}

func (m *UserModel) GetEntity() *entity.User {
	return m.entity
}

func (m *UserModel) GetEmptyEntity() entity.User {
	return entity.User{}
}

func (m *UserModel) IsEntityEmpty() bool {
	return *m.GetEntity() == m.GetEmptyEntity()
}

func (m *UserModel) IsEntityNil() bool {
	return m.GetEntity() == nil
}

func (m *UserModel) AddSupportsNum(supports_num uint64) {
	m.entity.SupportsRemainingNum += supports_num
	if m.entity.SupportsRemainingNum < 0 {
		m.entity.SupportsRemainingNum = 0
	}
	m.update_target_cols["supports_remaining_num"] = "supports_remaining_num"
}

func (m *UserModel) clear_update_target_cols() {
	m.update_target_cols = map[string]string{}
}

// DB Access //
// Fetch
func (m *UserModel) GetById(id int64) error {
	bean := m.GetEmptyEntity()
	_, err := m.BaseModel.GetById(id, &bean)
	if err != nil {
		return err
	}
	m.SetEntity(&bean)
	return nil
}

// Update
func (m *UserModel) UpdateById(_bean *entity.User, cols ...string) (int64, error) {
	if m.IsEntityNil() || m.IsEntityEmpty() {
		panic("Must set user entity at first")
	}

	rows, err := m.BaseModel.UpdateById(m.GetEntity().Id, _bean, cols...)
	if err != nil {
		return 0, err
	}
	if rows > 0 {
		m.SetEntity(_bean)
	}
	return rows, nil

}

// Insert
func (m *UserModel) Insert(_bean *entity.User) (int64, error) {
	m.BaseModel.SetNullableCols(_bean.GetInitNullableCols()...)
	rows, err := m.BaseModel.Insert(_bean, _bean.GetInitOmitCols()...)
	if err != nil {
		return 0, err
	}
	if rows > 0 {
		m.SetEntity(_bean)
	}
	return rows, nil
}

// ////////// //

// Interface //
