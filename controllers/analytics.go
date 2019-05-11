package controllers

import (
	"net/http"

	"github.com/philip-bui/fruits-service/errors"
	_ "github.com/philip-bui/fruits-service/protos"
	"github.com/philip-bui/fruits-service/services/auth"
	_ "github.com/philip-bui/fruits-service/services/postgres/dw"
)

// TODO:
func GetAnalytics(w http.ResponseWriter, req *http.Request) errors.HttpError {
	_, _ = auth.GetUserIDFromJWT(req.Header.Get(AuthorizationHeader))

	return nil
}
