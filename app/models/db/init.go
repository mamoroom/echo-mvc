package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"

	"github.com/mamoroom/echo-mvc/app/config"
	"github.com/mamoroom/echo-mvc/app/lib/logger"
	_ "github.com/mamoroom/echo-mvc/app/lib/util"

	_ "fmt"
	"os"
	"time"
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

	//cacher
	//cacher := xorm.NewLRUCacher2(xorm.NewMemoryStore(), time.Duration(conf.Db.Cache.ExpiresDuration)*time.Second, conf.Db.Cache.MaxElementSize)
	//DbMasterR.SetDefaultCacher(cacher)
}

func getDataSourceName(db_name string, db_connect_conf config.DbConnectConf) string {
	return db_connect_conf.UserName + ":" + db_connect_conf.Password + "@tcp(" + db_connect_conf.Host + ":" + db_connect_conf.Port + ")/" + db_name + "?charset=utf8"
}

func New(db_conf config.DbConf, db_connect_conf config.DbConnectConf) *xorm.Engine {
	var err error
	env := os.Getenv("CONFIGOR_ENV")

	db, err := xorm.NewEngine(db_conf.DriverName, getDataSourceName(db_conf.DbName, db_connect_conf))
	if err != nil {
		panic(err)
	}

	db.DatabaseTZ = time.Now().Location()
	db.TZLocation = time.Now().Location()

	// LOG
	conn_name := db_conf.DbName + "_" + db_connect_conf.Permission
	/*logfile, err := os.OpenFile(util.GetRootPath()+"/logs/"+conn_name+"_xorm.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("cannnot open xorm.log:" + err.Error())
	}*/
	logger := xorm.NewSimpleLogger(logger.GetRotateLogWriter("db/" + conn_name + "_xorm" + "_" + env))
	db.SetLogger(logger)

	if db_connect_conf.IsSqlLog {
		db.ShowSQL(true)
	}

	switch db_connect_conf.LogMode {
	case "DEBUG":
		db.Logger().SetLevel(core.LOG_DEBUG)
	case "INFO":
		db.Logger().SetLevel(core.LOG_INFO)
	case "WARN":
		db.Logger().SetLevel(core.LOG_WARNING)
	case "ERROR":
		db.Logger().SetLevel(core.LOG_ERR)
	}

	lifetime := time.Duration(db_connect_conf.ConnMaxLifetime) * time.Second
	db.SetMaxIdleConns(db_connect_conf.MaxIdleConns)
	db.SetMaxOpenConns(db_connect_conf.MaxOpenConns)
	db.DB().SetConnMaxLifetime(lifetime)

	//fmt.Printf(conn_name+"_xorm log level="+db_connect_conf.LogMode+" | idle: %d, open: %d, lifetime: %s", db_connect_conf.MaxIdleConns, db_connect_conf.MaxOpenConns, lifetime.String())

	return db
}
