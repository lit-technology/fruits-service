package fruits

import (
	"github.com/rs/zerolog/log"
)

var (
	InsertSurveyStmt = PrepareStatement(`
		INSERT INTO survey (user_id, name)
		VALUES ($1, $2)
		RETURNING id`)
	/*UpdateSurveyAnswersStmt = PrepareStatement(`
	UPDATE survey
	SET answers = answers + 1
	WHERE user_id = $1
	AND id = $2`)*/
)

func InsertSurvey(userID int64, name string) (int64, error) {
	var surveyID int64
	if err := InsertSurveyStmt.QueryRow(userID, name).Scan(&surveyID); err != nil {
		log.Error().Err(err).Int64("userID", userID).Str("name", name).Msg("error inserting survey")
		return 0, err
	}
	return surveyID, nil
}
