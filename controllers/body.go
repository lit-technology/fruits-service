package controllers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
)

const (
	AuthorizationHeader = "Authorization"
	ContentTypeHeader   = "Content-Type"
	ContentTypeJSON     = "application/json"
)

func ReadRequest(req *http.Request, pb proto.Unmarshaler) error {
	switch req.Header.Get(ContentTypeHeader) {
	case ContentTypeJSON:
		return ReadJSON(req.Body, pb)
	default:
		return ReadProtobuf(req.Body, pb)
	}
}

func ReadJSON(r io.Reader, pb proto.Unmarshaler) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		log.Error().Err(err).Msg("error reading bytes")
		return err
	}
	if err := json.Unmarshal(body, pb); err != nil {
		log.Error().Err(err).Msg("error unmarshalling bytes to protobuf json")
		return err
	}
	return nil
}

func ReadProtobuf(r io.Reader, pb proto.Unmarshaler) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		log.Error().Err(err).Msg("error reading bytes")
		return err
	}
	if err := pb.Unmarshal(body); err != nil {
		log.Error().Err(err).Interface("pb", pb).Msg("error unmarshalling bytes to protobuf")
		return err
	}
	return nil
}
