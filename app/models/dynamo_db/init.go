package dynamo_db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/mamoroom/echo-mvc/app/config"

	"os"
)

var conf = config.Conf
var cred *credentials.Credentials
var region *string

func init() {
	cred = credentials.NewStaticCredentials(conf.Aws.Credentials.AccessKeyId, conf.Aws.Credentials.SecretAccessKey, "") // 最後の引数は[セッショントークン]
	region = aws.String(getRegionFromInstanceMetaData())
}

func getRegionFromInstanceMetaData() (region string) {
	env := os.Getenv("CONFIGOR_ENV")
	if env == "local" || env == "dev" || env == "stg" {
		return conf.Aws.DynamoDb.Region
	}
	metadata := ec2metadata.New(session.New())
	region, err := metadata.Region()
	if err != nil {
		panic(err)
	}
	return region
}
