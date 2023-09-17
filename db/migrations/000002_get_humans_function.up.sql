create or replace function select_humans(
    pageNum int = 1,
    pageSize int = 10000
)
    returns table
        (
        id int,
        surname text,
        name text,
        patronymic text,
        age int,
        nationality text,
        gender text
        )
    LANGUAGE plpgsql
as
$$
begin
    return query
        SELECT *
        FROM humans h
        LIMIT pageSize
        OFFSET (pageNum - 1) * pageSize;
end
$$;