CREATE TABLE IF NOT EXISTS account (
    id SERIAL PRIMARY KEY,
    balance int NOT NULL default 0,
    created  date NOT NULL DEFAULT CURRENT_DATE
);

CREATE TABLE IF NOT EXISTS "user" (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    account_id int NOT NULL REFERENCES account (id),
    created  date NOT NULL DEFAULT CURRENT_DATE
    );