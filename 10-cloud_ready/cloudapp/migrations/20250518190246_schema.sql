-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

create table albums (
    id serial primary key,
    artist text,
    title text,
    year integer
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
