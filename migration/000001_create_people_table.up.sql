CREATE TABLE IF NOT EXISTS people (
    id serial PRIMARY KEY,
    name text NOT NULL,
    surname text NOT NULL,
    patronymic text,
    age int NOT NULL,
    gender text NOT NULL,
    nation text NOT NULL
)