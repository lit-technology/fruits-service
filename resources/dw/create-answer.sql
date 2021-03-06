CREATE TABLE IF NOT EXISTS user_agent (
	id INT GENERATED BY DEFAULT AS IDENTITY,
	agent VARCHAR(100)
);

CREATE INDEX IF NOT EXISTS user_agent_idx ON user_agent USING HASH(agent);

/*
CREATE UNLOGGED TABLE IF NOT EXISTS answer (
	survey_id BIGINT NOT NULL,
	id BIGINT PRIMARY KEY,
	user_id BIGINT,
	ip INET,
	user_agent INT,
	referrer VARCHAR(80),
	created TIMESTAMP NOT NULL DEFAULT NOW()
) PARTITION BY HASH(survey_id);

CREATE INDEX IF NOT EXISTS answer_idx ON answer USING BRIN(survey_id);

CREATE UNLOGGED TABLE IF NOT EXISTS answer_choice (
	survey_id BIGINT NOT NULL,
	answer_id BIGINT NOT NULL,
	question SMALLINT NOT NULL,
	choice SMALLINT,
	text VARCHAR(800)
) PARTITION BY HASH(survey_id);

CREATE INDEX IF NOT EXISTS answer_choice_idx ON answer_choice USING BRIN(survey_id);
*/
