package script

import (
	"github.com/mamoroom/echo-mvc/app/config"
	"github.com/Sirupsen/logrus"
)

var conf = config.Conf
var StatusSucceded = "succeded"
var StatusFailed = "failed"

//res object
type ResObject struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data interface{}

func resSucceded(data interface{}) *ResObject {
	return &ResObject{
		Status: StatusSucceded,
		Data:   data,
	}
}

func resFailed(data interface{}) *ResObject {
	return &ResObject{
		Status: StatusFailed,
		Data:   data,
	}
}

type ErrorParam struct {
	Logger    *logrus.Logger
	LogFunc   string
	ErrorType string
	Msg       string
	Param     map[string]interface{}
}

func handle_rollback_or_commit(f func() error) {
	err_tx := f()
	// for debug
	if err_tx != nil {
		panic(err_tx)
	}
}
