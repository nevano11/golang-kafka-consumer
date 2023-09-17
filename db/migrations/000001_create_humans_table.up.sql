CREATE TABLE IF NOT EXISTS humans
(
    id          serial PRIMARY KEY,
    surname     text    not null,
    name        text    not null,
    patronymic  text    null,
    age         integer not null,
    nationality text    not null,
    gender      text    not null
)