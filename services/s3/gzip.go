package s3

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"

	"github.com/rs/zerolog/log"
)

func Gzip(b []byte) ([]byte, error) {
	var buf bytes.Buffer
	// AWS Cloudfront uses 6 for maximum performance. Default 7.
	if w, err := gzip.NewWriterLevel(&buf, 7); err != nil {
		log.Error().Err(err).Msg("error creating gzip writer")
		return nil, err
	} else if _, err := w.Write(b); err != nil {
		log.Error().Err(err).Bytes("b", b).Msg("error gzipping bytes")
		return nil, err
	}
	return buf.Bytes(), nil
}

func GzipRead(body io.Reader) ([]byte, error) {
	r, err := gzip.NewReader(body)
	if err != nil {
		log.Error().Err(err).Msg("error creating gzip reader")
		return nil, err
	}
	defer r.Close()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Error().Err(err).Msg("error reading gzipped bytes")
		return nil, err
	}
	return b, nil
}
