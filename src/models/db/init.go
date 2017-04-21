package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"

	"github.com/mamoroom/echo-mvc/src/config"
	"github.com/mamoroom/echo-mvc/src/lib/util"

	"os"
)

var DbMasterR *xorm.Engine
var DbMasterW *xorm.Engine
var DbUserR *xorm.Engine
var DbUserW *xorm.Engine
var conf = config.Conf

func init() {
	DbMasterR = New(conf.Db.Master, conf.Db.Master.R)
	DbMasterW = New(conf.Db.Master, conf.Db.Master.W)
	DbUserR = New(conf.Db.User, conf.Db.User.R)
	DbUserW = New(conf.Db.User, conf.Db.User.W)
}

func getDataSourceName(db_name string, db_connect_conf config.DbConnectConf) string {
	return db_connect_conf.UserName + ":" + db_connect_conf.Password + "@tcp(" + db_connect_conf.Host + ":" + db_connect_conf.Port + ")/" + db_name + "?charset=utf8"
}

func New(db_conf config.DbConf, db_connect_conf config.DbConnectConf) *xorm.Engine {
	var err error
	db, err := xorm.NewEngine(db_conf.DriverName, getDataSourceName(db_conf.DbName, db_connect_conf))
	if err != nil {
		panic(err)
	}

	// LOG
	logfile, err := os.OpenFile(util.GetRootPath()+"/logs/xorm.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open xorm.log:" + err.Error())
	}
	logger := xorm.NewSimpleLogger(logfile)
	db.SetLogger(logger)

	if db_connect_conf.IsSqlLog {
		db.ShowSQL(true)
	}

	switch db_connect_conf.LogMode {
	case "debug":
		db.Logger().SetLevel(core.LOG_DEBUG)
	case "info":
		db.Logger().SetLevel(core.LOG_INFO)
	}

	return db
	//orm.SetMaxIdleConns(setting.MaxIdle)
	//orm.SetMaxOpenConns(setting.MaxOpen)
}
