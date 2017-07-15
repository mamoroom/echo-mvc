package dbh

import (
	"github.com/mamoroom/echo-mvc/app/models/db/db_conns"
	"github.com/go-xorm/xorm"

	"database/sql"
	_ "fmt"
	"reflect"
)

func New(db_conn db_conns.DbConn) *DbHandler {
	return &DbHandler{
		db_conn: db_conn,
	}
}

type DbHandler struct {
	db_conn db_conns.DbConn
	session *xorm.Session
}

// *xorm.Engine() or *xorm.Session()
// 必要なメソッドを都度追加する
type XormHandler interface {
	Id(id interface{}) *xorm.Session
	Where(query interface{}, args ...interface{}) *xorm.Session
	Insert(beans ...interface{}) (int64, error)
	Cols(columns ...string) *xorm.Session
	Omit(columns ...string) *xorm.Session
	Find(beans interface{}, condiBean ...interface{}) error
	Desc(columns ...string) *xorm.Session
	Asc(columns ...string) *xorm.Session
	NoCache() *xorm.Session
	//Exec(sqlStr string, args ...interface{}) (sql.Result, error)
}

/// Public Method ///
func (dbh *DbHandler) Handle() XormHandler {
	if dbh.session == nil {
		return dbh.get_engine()
	}
	return dbh.get_session()
}

func (dbh *DbHandler) GetDbConnName() string {
	return dbh.db_conn.GetDbConnName()
}

func (dbh *DbHandler) ClearCache(beans ...interface{}) error {
	return dbh.db_conn.ClearCache(beans...)
}

func (dbh *DbHandler) SetNewSession() {
	dbh.set_session(dbh.db_conn.NewSession())
}

func (dbh *DbHandler) BeginTx() error {
	err := dbh.session.Begin()
	if err != nil {
		return err
	}
	return nil
}

func (dbh *DbHandler) Rollback() error {
	return dbh.session.Rollback()
}

func (dbh *DbHandler) Commit() error {
	return dbh.session.Commit()
}

func (dbh *DbHandler) Exec(sqlStr string, args ...interface{}) (sql.Result, error) {
	return dbh.session.Exec(sqlStr, args...)
}

func (dbh *DbHandler) Close() {
	dbh.session.Close()
}

func (dbh *DbHandler) ForUpdate() {
	if dbh.session == nil {
		panic("Must start tx at first")
	}
	dbh.set_session(dbh.session.ForUpdate())
}

// [todo]: table依存のmethodがdbhにあるのはイケてないが...//
func (dbh *DbHandler) SetNullable(cols ...string) {
	s := dbh.get_session()
	dbh.set_session(s.Nullable(cols...))
}
func (dbh *DbHandler) GetTableName(bean interface{}) string {
	v := reflect.ValueOf(bean)
	return dbh.get_engine().TableMapper.Obj2Table(reflect.Indirect(v).Type().Name())
}

////////////////////////////////////////////////////////

/// Protected Method ///
func (dbh *DbHandler) get_engine() *xorm.Engine {
	return dbh.db_conn.Engine()
}

func (dbh *DbHandler) get_session() *xorm.Session {
	return dbh.session
}

func (dbh *DbHandler) set_session(s *xorm.Session) {
	dbh.session = s
}
