create function create_human(
    surname text,
    name text,
    patronymic text,
    age int,
    nationality text,
    gender text
) returns int
    LANGUAGE sql
as
$$
INSERT INTO humans (surname, name, patronymic, age, nationality, gender)
VALUES (surname, name, patronymic, age, nationality, gender)
RETURNING id;
$$;