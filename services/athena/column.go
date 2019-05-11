package athena

import (
	"fmt"
	"strconv"

	pb "github.com/philip-bui/fruits-service/protos"
)

type AthenaColumn struct {
	Name string
	Type string
}

func Max32(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func NewAthenaColumn(i int, q *pb.Question) *AthenaColumn {
	var colType string
	switch q.Type {
	case pb.Question_TEXT:
		colType = fmt.Sprintf("VARCHAR(%d)", Max32(q.MaxLength, 180))
	case pb.Question_PARAGRAPH:
		colType = fmt.Sprintf("VARCHAR(%d)", Max32(q.MaxLength, 800))
	case pb.Question_SINGLE, pb.Question_SLIDER, pb.Question_DROPDOWN:
		colType = "TINYINT" // 256
		break
	case pb.Question_MULTI, pb.Question_CHECKBOX, pb.Question_CROSSBOX:
		colType = "ARRAY<TINYINT>" // 256
		break
	case pb.Question_TABLE:
		colType = "MAP<TINYINT, TINYINT>" // 256
		break
	}
	return &AthenaColumn{Name: strconv.Itoa(i), Type: colType}
}
