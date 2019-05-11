package main

import (
	"bytes"
	"net/http"
	"testing"
	"time"

	pb "github.com/philip-bui/fruits-service/protos"
	"github.com/stretchr/testify/suite"
)

const (
	hostURL = "http://127.0.0.1:8080"
)

func TestMain(t *testing.M) {
	go func() {
		main()
	}()
	time.Sleep(1 * time.Second)
	t.Run()
}

func TestServer(t *testing.T) {
	suite.Run(t, new(MainSuite))
}

type MainSuite struct {
	suite.Suite
}

func (s *MainSuite) SetupSuite() {}

func (s *MainSuite) TestHealthCheck() {
	resp, err := http.Get(hostURL + "/health")
	s.NoError(err)
	s.NotNil(resp)
	s.Equal(http.StatusOK, resp.StatusCode)
}

func (s *MainSuite) TestPostSurvey() {
	reqBody := &pb.Survey{ID: 1, Name: "Test"}
	buf, err := reqBody.Marshal()
	s.NoError(err)

	resp, err := http.Post(hostURL+"/survey", ContentTypeProtobuf, bytes.NewBuffer(buf))
	s.NoError(err)
	s.NotNil(resp)
	s.Equal(http.StatusOK, resp.StatusCode)
}
