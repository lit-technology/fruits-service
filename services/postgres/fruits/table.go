package fruits

import (
	"strconv"
	"strings"

	pb "github.com/philip-bui/fruits-service/protos"
)

const (
	ColSmallInt      = "SMALLINT"
	ColSmallIntArray = "SMALLINT ARRAY"
)

func Max(x int32, y int64) int64 {
	if y >= int64(x) {
		return y
	}
	return int64(x)
}

func CreateTableForSurvey(s *pb.Survey) {
	sb := strings.Builder{}
	tableName := "s_" + strconv.FormatInt(s.ID, 10)
	sb.WriteString("CREATE TABLE " + tableName + "(")
	for i, q := range s.Questions {
		sb.WriteString(ColumnsForQuestion(strconv.FormatInt(int64(i), 64), q) + ",")
	}
	sb.WriteString("timestamp " +

		") PARTITION BY RANGE(timestamp); CREATE INDEX " + tableName + "_idx ON " + tableName + " USING BRIN(timestamp)")
}

func ColumnsForQuestion(name string, q *pb.Question) string {
	switch q.Type {
	case pb.Question_SINGLE, pb.Question_SLIDER, pb.Question_DROPDOWN:
		return name + " " + ColSmallInt
	case pb.Question_MULTI, pb.Question_CHECKBOX, pb.Question_CROSSBOX:
		return name + " " + ColSmallIntArray
	case pb.Question_TEXT, pb.Question_PARAGRAPH:
		return name + " " + ColVarchar(q)
	case pb.Question_TABLE:
		sb := strings.Builder{}
		for i := 0; i < len(q.Table.GetRows())-1; i++ {
			sb.WriteString(name + " " + ColSmallInt + ",")
		}
		sb.WriteString(name + " " + ColSmallInt)
		return sb.String()
	case pb.Question_TABLE_MULTI:
		sb := strings.Builder{}
		for i := 0; i < len(q.Table.GetRows())-1; i++ {
			sb.WriteString(name + " " + ColSmallIntArray + ",")
		}
		sb.WriteString(name + " " + ColSmallIntArray)
		return sb.String()
	default:
		panic("")
	}
	return ""
}

func ColVarchar(q *pb.Question) string {
	switch q.Type {
	case pb.Question_TEXT:
		return "VARCHAR(" + strconv.FormatInt(Max(q.MaxLength, 180), 64) + ")"
	case pb.Question_PARAGRAPH:
		return "VARCHAR(" + strconv.FormatInt(Max(q.MaxLength, 800), 64) + ")"
	default:
		panic("")
	}
}
