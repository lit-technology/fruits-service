# Surveys

## Question configuration

- Min length

Minimum length of choices, table column choices, or characters to be used. >= 1 means the question is required. Optional

- Max length

Maximum length of choices, table column choices, or characters to be used. Optional

## Survey serving

The survey forms is stored into S3, served through our CDN to deliver our forms closer to our users, whilst stored in web browsers cache.

Surveys are stored and indexed in SQL as follows. 

- UserID -> SurveyID.
- Survey Name -> SurveyID. `TODO:`

NOTE: We'll add updating survey question names or adding questions, but not updating the type or removing questions.

## Answers Dynamic SQL

Our questions have pre-defined answer bounds so naturally we store answers in a matching, optimal column type.

- Single Choice, Slider, Dropdown

Stored as `SMALLINT` representing the index of the choice made.

- Multi Choice, Checkbox, Crossbox

Stored as `Array<SMALLINT>` representing indices of the choices made.

- Text, Paragraph

Stored as `VARCHAR(180|800)` in UTF8.

- Table

Stored as `Array<SMALLINT>`. Array index represents row, value represents column choice.

- Table Multi

Stored as `Array<Array<<SMALLINT>>`. Array index represents row, value represents column choices.


- Additional columns are stored such as UserID, UserAgent, IP, Referrer.

### Further reading

https://developer.mozilla.org/en-US/docs/Web/HTTP/Browser_detection_using_the_user_agent
