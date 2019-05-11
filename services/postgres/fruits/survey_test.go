package fruits

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestServer(t *testing.T) {
	suite.Run(t, new(TransactionSuite))
}

type TransactionSuite struct {
	suite.Suite
	ctx    context.Context
	cancel context.CancelFunc
	tx     *sql.Tx
}

func (s *TransactionSuite) SetupTest() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), time.Duration(30000))
	tx, err := DB.BeginTx(s.ctx, &sql.TxOptions{})
	s.NoError(err)
	s.tx = tx
}

func (s *TransactionSuite) TearDownTest() {
	s.NoError(s.tx.Rollback())
}

func (s *TransactionSuite) TestInsertSurvey() {

}
