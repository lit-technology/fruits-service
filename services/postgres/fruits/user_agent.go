package fruits

var (
	InsertOrGetUserAgentIDStmt = PrepareStatement(`
		WITH row AS (
			INSERT INTO user_agent(agent)
			VALUES ($1)
			ON CONFLICT DO NOTHING
			RETURNING id
		)
		SELECT id FROM row
		UNION
		SELECT id FROM user_agent WHERE agent = $1
	`)
)

func InsertOrGetUserAgentID(userAgent string) (int32, error) {
	var ID int32
	if err := InsertOrGetUserAgentIDStmt.QueryRow(userAgent).Scan(&ID); err != nil {
		return 0, err
	}
	return ID, nil
}
