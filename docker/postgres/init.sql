CREATE DATABASE docker_test4;

CREATE SCHEMA IF NOT EXISTS sample AUTHORIZATION postgres;

\connect docker_test4

CREATE TABLE images (hash text, o_name text, n_name text, date date, size integer);