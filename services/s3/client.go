package s3

import (
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	BucketSurvey           = "fruits-surveys"
	ContentEncodingGzip    = "gzip"
	ContentTypeOctetStream = "application/octet-stream"
)

var (
	S3 *s3.S3
)

func init() {
	if sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"),
	}); err != nil {
		log.Fatal().Err(err).Msg("error connecting to s3")
	} else {
		S3 = s3.New(sess)
		log.Info().Str("endpoint", S3.Endpoint).Msg("connected to s3")
	}
}
