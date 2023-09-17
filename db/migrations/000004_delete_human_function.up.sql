create function delete_human(
    _id int
) returns int
    LANGUAGE sql
    as
$$
DELETE FROM humans
WHERE id = _id
RETURNING id;
$$;