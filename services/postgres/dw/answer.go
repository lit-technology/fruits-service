package dw

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/lib/pq"
	pb "github.com/philip-bui/fruits-service/protos"
	"github.com/philip-bui/fruits-service/services/postgres"
	"github.com/rs/zerolog/log"
)

var (
	AnswerPartitions = int64(2048)
)

var (
	ErrInvalidAnswer = errors.New("invalid answers to survey questions")
)

func ColumnString(i int) string {
	return "s_" + strconv.Itoa(i)
}

func ColumnInt(i int) string {
	return "n_" + strconv.Itoa(i)
}

func ColumnIntArray(i int) string {
	return "nn_" + strconv.Itoa(i)
}

func AnswerTable(surveyID int64) string {
	return "answer_" + fmt.Sprintf("%04d", surveyID%AnswerPartitions)
}

func ValidateString(s string, q *pb.Question) error {
	l := int32(len(s))
	if l < q.MinLength || (q.MaxLength != 0 && l > q.MaxLength) {
		return errors.New(fmt.Sprintf("text must contain between %d - %d characters", q.MinLength, q.MaxLength))
	}
	return nil
}

func ValidateChoice(c int32, q *pb.Question) error {
	if c >= int32(len(q.Choices)) || c < 0 {
		return errors.New(fmt.Sprintf("question has invalid choice %d from %d", c, len(q.Choices)))
	}
	return nil
}

func ValidateChoices(ints []int32, q *pb.Question) error {
	if len(ints) < int(q.MinLength) {
		return errors.New(fmt.Sprintf("question requires minimum %d choices", q.MinLength))
	} else if q.MaxLength > 0 && len(ints) > int(q.MaxLength) {
		return errors.New(fmt.Sprintf("question exceeds maximum %d choices", q.MaxLength))
	}
	for _, c := range ints {
		if c < 0 || int(c) >= len(q.Choices) {
			return errors.New(fmt.Sprintf("question has invalid choice %d from %d", c, len(q.Choices)))
		}
	}
	return nil
}

func ValidateTable(l []*pb.List, q *pb.Question) error {
	if len(l) != len(q.Table.Rows) {
		return errors.New("mismatched rows with answers")
	}
	for i, row := range l {
		if (row == nil || row.NS == nil) && q.MinLength == 0 {
			continue
		}

		if len(row.NS) < int(q.MinLength) {
			return errors.New(fmt.Sprintf("table row %d requires minimum %d columns", i, q.MinLength))
		} else if q.MaxLength > 0 && len(row.NS) > int(q.MaxLength) {
			return errors.New(fmt.Sprintf("table row %d exceeds maximum %d columns", i, q.MaxLength))
		}
		s := make(map[int32]interface{}, len(row.NS))
		for _, n := range row.NS {
			if n < 0 || int(n) >= len(q.Table.Columns) {
				return errors.New(fmt.Sprintf("table row %d selected invalid column %d from %d", i, n, len(q.Table.Columns)))
			}
			if _, ok := s[n]; ok {
				return errors.New(fmt.Sprintf("table row %d has duplicate column choice %d", i, n))
			}
			s[n] = nil
		}
	}
	return nil
}

func AnswerStringToSQL(a *pb.Answer) *sql.NullString {
	val := &sql.NullString{}
	if len(a.S) != 0 {
		val.Scan(a.S)
	}
	return val
}

func AnswerNumberToSQL(a *pb.Answer, q *pb.Question) *sql.NullInt64 {
	val := &sql.NullInt64{}
	if q.MinLength > 0 || !a.NULL {
		val.Scan(a.N)
	}
	return val
}

func AnswerTableToSQL(ls []*pb.List) ([]*sql.NullInt64, error) {
	if ls == nil {
		return nil, errors.New("error parsing answer list")
	}

	arr := make([]*sql.NullInt64, len(ls))
	for i, l := range ls {
		arr[i] = &sql.NullInt64{}
		if l == nil || l.NS == nil || len(l.NS) == 0 {
			continue
		}
		if err := arr[i].Scan(l.NS[0]); err != nil {
			log.Error().Err(err).Msg("error parsing list value to array")
			return nil, err
		}
	}
	return arr, nil
}

func AnswerTableMultiToSQL(ls []*pb.List, cols int) ([][]*sql.NullInt64, error) {
	if ls == nil {
		return nil, errors.New("error parsing answer list")
	}

	arr := make([][]*sql.NullInt64, len(ls))
	for i, l := range ls {
		row := make([]*sql.NullInt64, cols)
		currentCol := 0
		if l != nil && l.NS != nil {
			// We need to turn 3,1 to null,1,null,3.
			// Validate asserts no duplicate answers.
			SortInt32Array(l.NS)

			for j := 0; j < len(l.NS); j++ {
				col := int(l.NS[j])
				for currentCol < col {
					row[currentCol] = &sql.NullInt64{}
					currentCol++
				}

				v := &sql.NullInt64{}
				v.Scan(currentCol)
				row[currentCol] = v
				currentCol++
			}
		}
		for ; currentCol < cols; currentCol++ {
			row[currentCol] = &sql.NullInt64{}
		}
		arr[i] = row
	}
	return arr, nil
}

func InsertAnswerFromQuestionAnswer(surveyID, userID int64, ip, userAgent, referrer string, questions []*pb.Question, answers []*pb.Answer) error {
	if len(questions) != len(answers) {
		return ErrInvalidAnswer
	}
	columns := make([]string, len(questions))
	values := make([]interface{}, len(questions))
	for i, a := range answers {
		q := questions[i]
		switch q.Type {
		case pb.Question_TEXT, pb.Question_PARAGRAPH:
			if err := ValidateString(a.S, q); err != nil {
				log.Error().Err(err).Int("i", i).Msg("error validating question text")
				return ErrInvalidAnswer
			}
			columns[i] = ColumnString(i)
			values[i] = AnswerStringToSQL(a)

		case pb.Question_SINGLE, pb.Question_SLIDER, pb.Question_DROPDOWN:
			if err := ValidateChoice(a.N, q); err != nil {
				log.Error().Err(err).Int("i", i).Msg("error validating question single choice")
				return ErrInvalidAnswer
			}
			columns[i] = ColumnInt(i)
			values[i] = AnswerNumberToSQL(a, q)

		case pb.Question_MULTI, pb.Question_CHECKBOX, pb.Question_CROSSBOX:
			if err := ValidateChoices(a.GetNS(), q); err != nil {
				log.Error().Err(err).Int("i", i).Msg("error validating question multiple choices")
				return ErrInvalidAnswer
			}
			columns[i] = ColumnIntArray(i)
			SortInt32Array(a.NS)
			values[i] = pq.Array(a.NS)

		case pb.Question_TABLE, pb.Question_TABLE_MULTI:
			if err := ValidateTable(a.GetL(), q); err != nil {
				log.Error().Err(err).Int("i", i).Msg("error validating question table answers")
				return ErrInvalidAnswer
			}
			columns[i] = ColumnIntArray(i)

			switch q.Type {
			case pb.Question_TABLE:
				val, err := AnswerTableToSQL(a.L)
				if err != nil {
					log.Error().Err(err).Int("i", i).Msg("error parsing 1D array from table answer")
					return ErrInvalidAnswer
				}
				values[i] = pq.Array(val)
			case pb.Question_TABLE_MULTI:
				val, err := AnswerTableMultiToSQL(a.L, len(q.Table.Columns))
				if err != nil {
					log.Error().Err(err).Int("i", i).Msg("error parsing 2D array from table answer")
					return ErrInvalidAnswer
				}
				values[i] = pq.Array(val)
			}
		}
	}
	return InsertAnswer(surveyID, userID, ip, userAgent, referrer, columns, values)
}

func InsertAnswerQuery(surveyID int64, columns []string) string {
	return fmt.Sprintf(`INSERT INTO %s(%s, survey_id, user_id, ip, user_agent, referrer) VALUES ($1, $2, $3, $4, $5, %s)`,
		AnswerTable(surveyID), strings.Join(columns, ", "), postgres.NewVarArgs(6, 6+len(columns)))
}

func InsertAnswer(surveyID, userID int64, ip, userAgent, referrer string, columns []string, values []interface{}) error {
	values = append(values, surveyID)
	values = append(values, postgres.Int64Nullable(userID))
	values = append(values, postgres.StringNullable(ip))
	values = append(values, postgres.StringNullable(userAgent))
	values = append(values, postgres.StringNullable(referrer))
	log.Debug().Interface("columns", columns).Interface("values", values).Msg("error inserting answer")
	if _, err := DB.Exec(InsertAnswerQuery(surveyID, columns), values...); err != nil {
		log.Error().Err(err).Interface("columns", columns).Interface("values", values).Msg("error inserting answer")
		return err
	}
	return nil
}
