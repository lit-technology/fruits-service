package athena

import (
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
)

const (
	BucketAnswer = "fruits-answers"
)

var (
	Athena *athena.Athena
)

func init() {
	if sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"),
	}); err != nil {
		log.Fatal().Err(err).Msg("error connecting to Athena")
	} else {
		Athena = athena.New(sess)
		log.Info().Str("endpoint", Athena.Endpoint).Msg("connected to Athena")
	}
}
