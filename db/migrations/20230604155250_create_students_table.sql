-- migrate:up

create table if not exists students
(
    id         varchar not null primary key,
    name       varchar not null,
    cpf        varchar not null,
    birth_date date    not null,
    email      varchar not null,
    secret     varchar not null
);

-- migrate:down
drop table if exists students

