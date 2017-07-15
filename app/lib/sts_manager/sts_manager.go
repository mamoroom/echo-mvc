package sts_manager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/mamoroom/echo-mvc/app/lib/util"

	"time"
)

type StsManager struct {
	*sts.STS
	user_id      int64
	nazo_id      int64
	file_name    string
	content_type string
	time         time.Time
}

type ResUrls struct {
	S3PutUrl string `json:"s3_put_url"`
	MediaUrl string `json:"media_url"`
}

func New(user_id int64, nazo_id int64, file_name string, content_type string, t time.Time) *StsManager {
	_sts := sts.New(session.New(&aws.Config{
		Region:      region,
		Credentials: cred,
	}))
	return &StsManager{
		STS:          _sts,
		user_id:      user_id,
		nazo_id:      nazo_id,
		file_name:    file_name,
		content_type: content_type,
		time:         t,
	}
}

func (s *StsManager) GetUrls() (*ResUrls, error) {
	params := &sts.AssumeRoleInput{
		RoleArn:         aws.String(conf.Aws.Sts.RoleArn),         // Required
		RoleSessionName: aws.String(conf.Aws.Sts.RoleSessionName), // Required
		/*DurationSeconds: 120,
		  ExternalId:      aws.String("externalIdType"),
		  Policy:          aws.String("sessionPolicyDocumentType"),
		  SerialNumber:    aws.String("serialNumberType"),
		  TokenCode:       aws.String("tokenCodeType"),*/
	}
	resp, err := s.AssumeRole(params)

	if err != nil {
		return nil, err
	}

	svc := s3.New(session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(*resp.Credentials.AccessKeyId, *resp.Credentials.SecretAccessKey, *resp.Credentials.SessionToken),
		Region:      region,
	}))

	file_name_suffix := util.GetFileNameSuffix(s.file_name)
	key := conf.Aws.Sts.S3SignedUrl.KeyPrefix + "/" + util.CastInt64ToStr(s.nazo_id) + "/" + util.CastInt64ToStr(s.user_id) + "/" + util.GetTimestampWithRand8(s.time) + "." + file_name_suffix
	expires_sec := time.Duration(conf.Aws.Sts.S3SignedUrl.Expires) * time.Second
	//expires_time := s.time.Add(expires_sec)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(conf.Aws.Sts.S3SignedUrl.BucketName),
		Key:         aws.String(key),
		ContentType: aws.String(s.content_type),
		// [todo]: バグるのでいったんコメントアウト
		//Expires:     &expires_time,
	})
	s3_put_url, err := req.Presign(expires_sec)

	if err != nil {
		return nil, err
	}

	return &ResUrls{
		S3PutUrl: s3_put_url,
		//MediaUrl: "https://" + conf.Aws.Sts.S3SignedUrl.BucketName + "." + conf.Server.UserMediaDomain + "/" + key,
		MediaUrl: conf.Server.UserMediaDomain + "/" + conf.Aws.Sts.S3SignedUrl.BucketName + "/" + key,
	}, err
}
