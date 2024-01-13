CREATE USER test_user WITH PASSWORD 'test_password';

CREATE DATABASE service WITH OWNER = test_user;

\c service test_user;