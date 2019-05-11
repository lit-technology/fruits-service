package athena

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/athena"
)

func surveyAnswerLocation(userID, surveyID int64) string {
	return fmt.Sprintf("s3://%s/%d/%d", BucketAnswer, userID, surveyID)
}

// https://docs.aws.amazon.com/athena/latest/ug/create-table.html
func CreateTable(tableName, tableColumns string, userID, surveyID int64) {
	Athena.StartQueryExecution(&athena.StartQueryExecutionInput{
		QueryString: aws.String(fmt.Sprintf(`
			CREATE EXTERNAL TABLE %s (
				user_id BIGINT,
				ip_address VARCHAR(39),
				%s
			)
			PARTITIONED BY(year int, month int, day int)
			STORED AS PARQUET
			LOCATION s3://%s/%d/%d
			tblproperties ("parquet.compress"="SNAPPY")
		`, tableName, tableColumns, BucketAnswer, userID, surveyID)),
		ResultConfiguration: &athena.ResultConfiguration{
			OutputLocation: aws.String("s3://" + BucketAnswer),
		},
	})
}

func QueryTable() {
	// TODO: Get chart format and query args format.
}
