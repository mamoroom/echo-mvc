package kinesis

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	SDK "github.com/aws/aws-sdk-go/service/kinesis"
)

type Kinesis struct {
	client *SDK.Kinesis
	stream *Stream
}

type Stream struct {
	name string
}

func New() *Kinesis {
	_client := SDK.New(session.New(&aws.Config{
		Credentials: cred,
		Region:      region,
	}))

	return &Kinesis{
		client: _client,
		stream: &Stream{
			conf.Aws.Kinesis.Stream.Name,
		},
	}
}
