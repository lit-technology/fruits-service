package controllers

import (
	"net/http"

	"github.com/philip-bui/fruits-service/errors"
	pb "github.com/philip-bui/fruits-service/protos"
	"github.com/philip-bui/fruits-service/services/auth"
	"github.com/philip-bui/fruits-service/services/postgres/dw"
	"github.com/philip-bui/fruits-service/services/s3"
)

func PostAnswer(w http.ResponseWriter, req *http.Request) errors.HttpError {
	a := &pb.SurveyAnswer{}
	if err := ReadRequest(req, a); err != nil {
		return errors.ErrBadRequest
	}
	userID, _ := auth.GetUserIDFromJWT(req.Header.Get(AuthorizationHeader))
	s, err := s3.GetSurvey(a.SurveyUserID, a.SurveyID)
	if err != nil {
		return errors.ErrNotFound
	}
	if err := dw.InsertAnswerFromQuestionAnswer(a.SurveyID, userID,
		req.RemoteAddr, req.UserAgent(), a.Referrer, s.Questions, a.Answers); err != nil {
		if err == dw.ErrInvalidAnswer {
			return errors.ErrBadRequest
		} else {
			return errors.ErrInternalServer
		}
	}
	return nil
}
