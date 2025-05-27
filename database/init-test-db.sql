CREATE DATABASE bd2_test;

\c bd2_test;

\i /docker-entrypoint-initdb.d/01-schema.sql

INSERT INTO manufacturers (name) VALUES ('BMW'), ('Audi'), ('Toyota'), ('Honda'), ('Aston Martin');

SELECT 'Test database bd2_test created the schema successfully' AS status;
