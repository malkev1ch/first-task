CREATE TABLE cats
(
    row_id     SERIAL    NOT NULL UNIQUE,
    name       VARCHAR   NOT NULL,
    date_birth DATE,
    vaccinated boolean,
    created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc')
)