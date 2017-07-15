package sts_manager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/mamoroom/echo-mvc/app/config"
)

var conf = config.Conf
var cred *credentials.Credentials
var region *string

func init() {
	cred = credentials.NewStaticCredentials(conf.Aws.Credentials.AccessKeyId, conf.Aws.Credentials.SecretAccessKey, "") // 最後の引数は[セッショントークン]
	region = aws.String(conf.Aws.Sts.Region)
}
