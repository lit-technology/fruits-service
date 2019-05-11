package dw

import (
	"database/sql"
	"math/rand"
	"testing"

	pb "github.com/philip-bui/fruits-service/protos"
	"github.com/philip-bui/fruits-service/services/postgres"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestAnswerTable(t *testing.T) {
	assert.Equal(t, "answer_0000", AnswerTable(AnswerPartitions))
	assert.Equal(t, "answer_0001", AnswerTable(AnswerPartitions+1))
}

func TestAnswerStringToSQL(t *testing.T) {
	assert.NoError(t, ValidateString("", &pb.Question{}))

	assert.Error(t, ValidateString("", &pb.Question{MinLength: 1}))
	assert.False(t, AnswerStringToSQL(&pb.Answer{S: ""}).Valid)

	assert.NoError(t, ValidateString("", &pb.Question{MaxLength: 1}))
	assert.NoError(t, ValidateString("a", &pb.Question{MinLength: 1}))
	assert.True(t, AnswerStringToSQL(&pb.Answer{S: "a"}).Valid)
}

func TestAnswerNumberToSQL(t *testing.T) {
	assert.NoError(t, ValidateChoice(0, &pb.Question{Choices: NewQuestionChoices(1)}))

	assert.Error(t, ValidateChoice(-1, &pb.Question{Choices: NewQuestionChoices(1)}))
	assert.Error(t, ValidateChoice(2, &pb.Question{Choices: NewQuestionChoices(2)}))

	assert.True(t, AnswerNumberToSQL(&pb.Answer{N: 0}, &pb.Question{Choices: NewQuestionChoices(1)}).Valid)
	assert.False(t, AnswerNumberToSQL(&pb.Answer{NULL: true}, &pb.Question{Choices: NewQuestionChoices(1)}).Valid)
	assert.False(t, AnswerNumberToSQL(&pb.Answer{N: 0, NULL: true}, &pb.Question{Choices: NewQuestionChoices(1)}).Valid)
}

func NewListAnswer(values ...interface{}) *pb.Answer {
	ans := &pb.Answer{L: make([]*pb.List, len(values))}
	for i, v := range values {
		if v == nil {
			continue
		}
		if l, ok := v.([]int32); ok {
			ans.L[i] = &pb.List{
				NS: l,
			}
		}
	}
	return ans
}

func TestAnswerTableToSQL(t *testing.T) {
	ans := NewListAnswer([]int32{0}, nil, []int32{3})
	q := &pb.Question{
		Table: NewTable(3, 4),
	}
	assert.NoError(t, ValidateTable(ans.L, q))

	arr, err := AnswerTableToSQL(ans.L)
	assert.NoError(t, err)
	assert.Equal(t, []*sql.NullInt64{postgres.Int64(0), &postgres.NullInt64, postgres.Int64(3)}, arr)
	assert.NotEqual(t, []*sql.NullInt64{}, arr)
	assert.NotEqual(t, []*sql.NullInt64{postgres.Int64(0)}, arr)
	assert.NotEqual(t, []*sql.NullInt64{postgres.Int64(0), &postgres.NullInt64, postgres.Int64(3), &postgres.NullInt64}, arr)

	assert.Error(t, ValidateTable(ans.L, &pb.Question{
		Table: NewTable(2, 5),
	}))
	assert.Error(t, ValidateTable(ans.L, &pb.Question{
		Table: NewTable(3, 3),
	}))
	assert.Error(t, ValidateTable(ans.L, &pb.Question{
		Table: NewTable(2, 2),
	}))
}

func TestAnswerTableMultiToSQL(t *testing.T) {
	ans := NewListAnswer([]int32{0}, nil, []int32{3, 1})
	q := &pb.Question{
		Table: NewTable(3, 5),
	}
	assert.NoError(t, ValidateTable(ans.L, q))
	arr, err := AnswerTableMultiToSQL(ans.L, 5)
	assert.NoError(t, err)
	assert.Equal(t, len(ans.L), len(arr))
	for _, arr2 := range arr {
		assert.Equal(t, 5, len(arr2))
	}
	// Assert 0, nil, nil, nil, nil
	for i, v := range arr[0] {
		if i == 0 {
			assert.Equal(t, postgres.Int64(0), v, "Expected row 0 column 0 to have valid 0")
			continue
		}
		assert.Equal(t, &postgres.NullInt64, v, "Expected row 0 to have nils")
	}
	// Assert nil, nil, nil, nil, nil
	for _, v := range arr[1] {
		assert.Equal(t, &postgres.NullInt64, v, "Expected row 1 to have nils")
	}
	// Assert nil, nil, nil, 3, nil
	for i, v := range arr[2] {
		if i == 1 {
			assert.Equal(t, postgres.Int64(1), v, "Expected row 2 column 1 to have valid 1")
			continue
		} else if i == 3 {
			assert.Equal(t, postgres.Int64(3), v, "Expected row 2 column 3 to have valid 3")
			continue
		}
		assert.Equal(t, &postgres.NullInt64, v, "Expected row 2 to have nils")
	}

	zerolog.SetGlobalLevel(zerolog.FatalLevel)

	// NOTE: Above answer is 3 rows, 4 columns.
	assert.Error(t, ValidateTable(ans.L, &pb.Question{Table: NewTable(3, 3)}))
	assert.Error(t, ValidateTable(ans.L, &pb.Question{Table: NewTable(2, 5)}))

	ans.L[0].NS[0] = -1
	assert.Error(t, ValidateTable(ans.L, q))
}

func TestInsertAnswerStmt(t *testing.T) {
	s := InsertAnswerQuery(2048, []string{"S_1", "N_2", "NN_3"})
	assert.Equal(t, "INSERT INTO answer_0000(S_1, N_2, NN_3, survey_id, user_id, ip, user_agent, referrer) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", s)
	s = InsertAnswerQuery(1, []string{"S_1"})
	assert.Equal(t, "INSERT INTO answer_0001(S_1, survey_id, user_id, ip, user_agent, referrer) VALUES ($1, $2, $3, $4, $5, $6)", s)
}

func NewQuestionChoices(choices int) []*pb.QuestionChoice {
	return make([]*pb.QuestionChoice, choices)
}

func NewTable(rows, cols int) *pb.Table {
	t := &pb.Table{}
	t.Rows = make([]*pb.Row, rows)
	t.Columns = make([]*pb.Column, cols)
	return t
}

func NewQuestionAndAnswer() (*pb.Question, *pb.Answer) {
	t := rand.Int31() % 10
	q := &pb.Question{}
	a := &pb.Answer{}
	switch t {
	case 0:
		q.Type = pb.Question_SINGLE
		q.Choices = NewQuestionChoices(7)
		a.N = 3
	case 1:
		q.Type = pb.Question_MULTI
		q.Choices = NewQuestionChoices(7)
	case 2:
		q.Type = pb.Question_TEXT
		a.S = "Test"
	case 3:
		q.Type = pb.Question_PARAGRAPH
	case 4:
		q.Type = pb.Question_SLIDER
		q.Choices = NewQuestionChoices(7)
		a.NS = []int32{0}
	case 5:
		q.Type = pb.Question_CHECKBOX
		q.Choices = NewQuestionChoices(7)
		a.NS = []int32{6}
	case 6:
		q.Type = pb.Question_CROSSBOX
		q.Choices = NewQuestionChoices(7)
	case 7:
		q.Type = pb.Question_DROPDOWN
		q.Choices = NewQuestionChoices(7)
	case 8:
		q.Type = pb.Question_TABLE
		q.Table = NewTable(3, 5)
		a.L = make([]*pb.List, len(q.Table.Rows))
	default:
		q.Type = pb.Question_TABLE_MULTI
		q.Table = NewTable(3, 5)
		a.L = make([]*pb.List, len(q.Table.Rows))
	}
	return q, a
}

func BenchmarkInsertAnswerParallel(b *testing.B) {
	b.RunParallel(func(bp *testing.PB) {
		for bp.Next() {
			surveyID := rand.Int63()
			n := 10
			questions := make([]*pb.Question, n)
			answers := make([]*pb.Answer, n)
			for i := 0; i < n; i++ {
				questions[i], answers[i] = NewQuestionAndAnswer()
			}
			assert.NoError(b, InsertAnswerFromQuestionAnswer(surveyID, 0, "", "answer_test.go", "philip", questions, answers))
		}
	})
}
