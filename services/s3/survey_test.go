package s3

import (
	"testing"

	pb "github.com/philip-bui/fruits-service/protos"
	"github.com/stretchr/testify/suite"
)

func TestServer(t *testing.T) {
	suite.Run(t, new(LRUCacheSuite))
}

type LRUCacheSuite struct {
	suite.Suite
}

func (s *LRUCacheSuite) BeforeTest() {
	Cache.Purge()
}

func (s *LRUCacheSuite) TestLRUCache() {
	s0 := &pb.Survey{
		ID: 0,
	}
	// Test cache add and get
	Cache.Add(0, s0)
	val, ok := Cache.Get(0)
	s.True(ok)
	s.Equal(s0, val)

	// Test cache item after removal
	Cache.Remove(0)
	s.Equal(s0, val)

	// Test cache after remove
	val, ok = Cache.Get(0)
	s.False(ok)
	s.Nil(val)
}
