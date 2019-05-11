package controllers

import (
	"net/http"

	"github.com/philip-bui/fruits-service/errors"
	pb "github.com/philip-bui/fruits-service/protos"
	"github.com/philip-bui/fruits-service/services/auth"
	"github.com/philip-bui/fruits-service/services/postgres/fruits"
	"github.com/philip-bui/fruits-service/services/s3"
)

func PostSurvey(w http.ResponseWriter, req *http.Request) errors.HttpError {
	s := &pb.Survey{}
	userID, err := auth.GetUserIDFromJWT(req.Header.Get(AuthorizationHeader))
	if err := ReadRequest(req, s); err != nil {
		return errors.ErrBadRequest
	}
	if err := ValidateSurvey(s); err != nil {
		return errors.ErrBadRequest
	}
	surveyID, err := fruits.InsertSurvey(userID, s.Name)
	if err != nil {
		return errors.ErrInternalServer
	}
	if err := s3.UploadSurvey(userID, surveyID, s); err != nil {
		return errors.ErrInternalServer
	}
	return nil
}
