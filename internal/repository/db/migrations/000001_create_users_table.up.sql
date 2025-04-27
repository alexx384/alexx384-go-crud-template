CREATE TABLE IF NOT EXISTS users (
    id int primary key generated always as identity,
    name VARCHAR(255) not null,
    email VARCHAR(255) not null,
    age INT not null
);
