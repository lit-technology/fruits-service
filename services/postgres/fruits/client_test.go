package fruits

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// This tests the DB connection, and prepared statements, ensuring the statements are valid SQL syntax.
// For non-prepared statements, please add test to ensure it runs.
func TestDB(t *testing.T) {
	assert.NotNil(t, DB)
}
