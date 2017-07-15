package logger

import (
	"github.com/mamoroom/echo-mvc/app/lib/util"
	"github.com/Sirupsen/logrus"
	"github.com/doloopwhile/logrusltsv"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"

	"os"
	"time"
)

var ErrorLogger *logrus.Logger
var AppLogger *logrus.Logger
var PointsRankingBatchLogger *logrus.Logger
var NazoRankingBatchLogger *logrus.Logger
var SendMailBatchLogger *logrus.Logger
var SystemLogger *logrus.Logger

func init() {
	env := os.Getenv("CONFIGOR_ENV")
	ErrorLogger = New("error/error"+"_"+env, &logrusltsv.Formatter{}, logrus.WarnLevel)
	AppLogger = New("app/app"+"_"+env, &logrus.JSONFormatter{}, logrus.WarnLevel)
	PointsRankingBatchLogger = New("batch/points_rank/log"+"_"+env, &logrusltsv.Formatter{}, logrus.InfoLevel)
	NazoRankingBatchLogger = New("batch/nazo_rank/log"+"_"+env, &logrusltsv.Formatter{}, logrus.InfoLevel)
	SendMailBatchLogger = New("batch/send_mail/log"+"_"+env, &logrusltsv.Formatter{}, logrus.InfoLevel)
	SystemLogger = New("system/log"+"_"+env, &logrus.JSONFormatter{}, logrus.InfoLevel)
}

func New(path string, formatter logrus.Formatter, log_level logrus.Level) *logrus.Logger {
	var log = logrus.New()
	log.Out = GetRotateLogWriter(path)
	log.Formatter = formatter
	log.Level = log_level
	return log
}

func GetRotateLogWriter(path string) *rotatelogs.RotateLogs {
	log_file_path := util.GetRootPath() + "/../logs/" + path

	logf, err := rotatelogs.New(
		log_file_path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(log_file_path),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic("failed to create rotatelogs: " + err.Error())
	}
	return logf
}
