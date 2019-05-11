// fruits-service project fruits-service.go
package main

import (
	"net/http"
	"os"
	"time"

	"github.com/philip-bui/fruits-service/controllers"
	"github.com/philip-bui/fruits-service/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	ContentTypeProtobuf = "application/protobuf"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC1123
	switch os.Getenv("LOG") {
	case "ERROR":
		log.Logger = log.Level(zerolog.ErrorLevel)
	case "WARN":
		log.Logger = log.Level(zerolog.WarnLevel)
	case "INFO":
		log.Logger = log.Level(zerolog.InfoLevel)
	default:
		log.Logger = log.Level(zerolog.DebugLevel)
	}
	if log.Debug().Enabled() {
		log.Logger = log.With().Caller().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

type Middleware = func(http.ResponseWriter, *http.Request) errors.HttpError

func main() {
	HandlePOST("/survey", controllers.PostSurvey)
	http.HandleFunc("/health", HealthCheck)
	log.Info().Msg("listening to :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal().Err(err).Msg("error listening to 8080")
	}
}

func HealthCheck(http.ResponseWriter, *http.Request) {

}

func HandleGET(pattern string, handler Middleware) {
	HandleMethod(pattern, http.MethodGet, handler)
}

func HandlePOST(pattern string, handler Middleware) {
	HandleMethod(pattern, http.MethodPost, handler)
}

func HandleMethod(pattern, method string, handler Middleware) {
	http.HandleFunc(pattern, ErrorMiddleware(LogMiddleware(handler), method))
}

func ErrorMiddleware(handler Middleware, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err errors.HttpError
		if req.Method != method {
			err = errors.ErrMethodNotAllowed
		} else {
			err = handler(w, req)
		}
		if err != nil {
			http.Error(w, err.Error(), err.Code())
		}
	}
}

func LogMiddleware(handler Middleware) Middleware {
	return func(w http.ResponseWriter, req *http.Request) errors.HttpError {
		start := time.Now()
		log := func(level zerolog.Level) *zerolog.Event {
			return log.WithLevel(level).Str("method", req.Method).
				Str("path", req.URL.String()).
				Str("ip", req.RemoteAddr).
				Str("user-agent", req.UserAgent()).
				Dur("duration", time.Since(start))
		}
		if err := handler(w, req); err != nil {
			log(zerolog.ErrorLevel).Int("status", err.Code()).Msg(err.Error())
			return err
		} else {
			log(zerolog.DebugLevel).Msg(err.Error())
		}
		return nil
	}
}
