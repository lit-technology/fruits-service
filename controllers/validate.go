package controllers

import (
	"errors"
	"fmt"
	"strings"

	pb "github.com/philip-bui/fruits-service/protos"
)

func ValidateSurvey(s *pb.Survey) error {
	s.Name = strings.TrimSpace(s.Name)
	if len(s.Name) == 0 {
		return errors.New("surveys must contain a name")
	}
	if s.Questions == nil || len(s.Questions) == 0 {
		return errors.New("surveys must contain questions")
	}
	return nil
}

func ValidateSurveyQuestion(q *pb.Question) error {
	q.Name = strings.TrimSpace(q.Name)
	if len(q.Name) == 0 {
		return errors.New("questions must contain a name")
	}
	if q.MaxLength != 0 && q.MaxLength < q.MinLength {
		return errors.New(fmt.Sprintf("question maximum length is less than minimum length"))
	}
	switch q.Type {
	case pb.Question_TEXT, pb.Question_PARAGRAPH:
		break
	case pb.Question_TABLE, pb.Question_TABLE_MULTI:
		if q.Table == nil || q.Table.Rows == nil || len(q.Table.Rows) == 0 {
			return errors.New("table questions must contain rows")
		} else if q.Table.Columns == nil || len(q.Table.Columns) == 0 {
			return errors.New("table questions must contain columns")
		} else if int(q.MinLength) > len(q.Table.Columns) {
			return errors.New("minimum choices exceeds columns")
		}
		for i, row := range q.Table.Rows {
			row.Name = strings.TrimSpace(row.Name)
			if len(row.Name) == 0 {
				return errors.New(fmt.Sprintf("row %d is missing a name", i))
			}
		}
		for i, col := range q.Table.Columns {
			col.Name = strings.TrimSpace(col.Name)
			if len(col.Name) == 0 {
				return errors.New(fmt.Sprintf("column %d is missing a name", i))
			}
		}
	default:
		if int(q.MinLength) > len(q.Choices) {
			return errors.New("minimum choices exceeds choices")
		}
	}
	return nil
}
