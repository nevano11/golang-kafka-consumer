create function edit_human(
    _id int,
    _surname text,
    _name text,
    _patronymic text,
    _age int,
    _nationality text,
    _gender text
) returns int
    LANGUAGE sql
    as
$$
UPDATE humans h
SET
    surname = _surname,
    name = _name,
    patronymic = _patronymic,
    age = _age,
    nationality = _nationality,
    gender = _gender
WHERE id = _id
RETURNING id;
$$;