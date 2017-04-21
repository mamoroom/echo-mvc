package dbh

import (
	"github.com/go-xorm/xorm"
	"github.com/mamoroom/echo-mvc/src/models/db/db_conns"

	_ "fmt"
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

func (dbh *DbHandler) SetNewSession() {
	dbh.set_session(dbh.db_conn.NewSession())
}

func (dbh *DbHandler) SetNullable(cols []string) {
	s := dbh.get_session()
	dbh.set_session(s.Nullable(cols...))
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

func (dbh *DbHandler) Close() {
	dbh.session.Close()
}

func (dbh *DbHandler) ForUpdate() {
	if dbh.session == nil {
		panic("Must start tx at first")
	}
	dbh.set_session(dbh.session.ForUpdate())
}

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
