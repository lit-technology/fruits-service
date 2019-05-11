CREATE TABLE IF NOT EXISTS survey (
	user_id BIGINT NOT NULL,
	id BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
	name VARCHAR(100),
	index SMALLINT,
	created TIMESTAMP NOT NULL DEFAULT NOW(),
	deleted TIMESTAMP
);

CREATE INDEX IF NOT EXISTS survey_idx ON survey USING BTREE(user_id, index);
