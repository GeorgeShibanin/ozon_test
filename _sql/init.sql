CREATE TABLE [IF NOT EXISTS] postgres (
     id varchar(10) UNIQUE PRIMARY KEY,
     url varchar UNIQUE NOT NULL
);