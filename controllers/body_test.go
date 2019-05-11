package controllers

import (
	"bytes"
	"testing"

	pb "github.com/philip-bui/fruits-service/protos"
	"github.com/stretchr/testify/assert"
)

func TestReadJSON(t *testing.T) {
	s := &pb.Survey{}
	ReadJSON(bytes.NewBuffer([]byte(`{"ID": 1}`)), s)
	assert.Equal(t, int64(1), s.ID)
}

func TestReadProtobuf(t *testing.T) {
	s := &pb.Survey{ID: 1, Name: "Test"}
	b, err := s.Marshal()
	assert.NoError(t, err)
	assert.NotEmpty(t, b)

	s2 := &pb.Survey{}
	ReadProtobuf(bytes.NewBuffer(b), s2)
	assert.Equal(t, s, s2)
	assert.Equal(t, s.ID, s2.ID)
	assert.Equal(t, s.Name, s2.Name)
}
