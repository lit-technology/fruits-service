package s3

import (
	"bytes"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/golang-lru"
	pb "github.com/philip-bui/fruits-service/protos"
	"github.com/rs/zerolog/log"
)

var (
	Cache *lru.Cache
)

func init() {
	cache, err := lru.New(1024)
	if err != nil {
		log.Fatal().Err(err).Int("size", 1024).Msg("error creating lru cache")
	}
	Cache = cache
}

func surveyKey(userID, surveyID int64) *string {
	// We use Base 10 for readability, but upgrade to Base 2^ for performance later on
	return aws.String(strconv.FormatInt(userID, 10) + "/" + strconv.FormatInt(surveyID, 10))
}

func GzipSurvey(s *pb.Survey) ([]byte, error) {
	b, err := s.Marshal()
	if err != nil {
		log.Error().Err(err).Msg("error marshalling survey to bytes")
		return nil, err
	}
	return Gzip(b)
}

func UploadSurvey(userID, surveyID int64, s *pb.Survey) error {
	b, err := GzipSurvey(s)
	if err != nil {
		log.Error().Err(err).Msg("error gzipping survey")
		return err
	}
	if _, err := S3.PutObject(&s3.PutObjectInput{
		Body:            bytes.NewReader(b),
		Bucket:          aws.String(BucketSurvey),
		Key:             surveyKey(userID, surveyID),
		ContentEncoding: aws.String(ContentEncodingGzip),
		ContentLength:   aws.Int64(int64(len(b))),
		ContentType:     aws.String(ContentTypeOctetStream),
	}); err != nil {
		log.Error().Err(err).Msg("error uploading survey to s3")
		return err
	}
	return nil
}

func GetSurvey(userID int64, surveyID int64) (*pb.Survey, error) {
	if v, ok := Cache.Get(surveyID); ok {
		if s, ok := v.(*pb.Survey); ok {
			return s, nil
		}
		log.Error().Int64("surveyID", surveyID).Interface("interface", v).Msg("found unexpected object for survey cache")
	}
	obj, err := S3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(BucketSurvey),
		Key:    surveyKey(userID, surveyID),
	})
	if err != nil {
		log.Error().Err(err).Msg("error getting survey from s3")
		return nil, err
	}
	defer obj.Body.Close()
	b, err := GzipRead(obj.Body)
	if err != nil {
		log.Error().Err(err).Msg("error reading gzipped survey from s3")
		return nil, err
	}
	s := &pb.Survey{}
	if err := s.Unmarshal(b); err != nil {
		log.Error().Err(err).Msg("error un-marshalling bytes to survey")
		return nil, err
	}
	Cache.Add(surveyID, s)
	return s, nil
}
