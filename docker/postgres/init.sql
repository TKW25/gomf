CREATE DATABASE image_database;

\connect image_database

CREATE SCHEMA IF NOT EXISTS sample AUTHORIZATION postgres;

CREATE TABLE images (hash text, o_name text, n_name text, date date, size integer);

INSERT INTO images ("hash", "o_name", "n_name", "date", "size")
       VALUES('sample_hash', 'fake_file', 'fake_file', '2019-01-01', 0);
