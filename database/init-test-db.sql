CREATE DATABASE bd2_test;

\c bd2_test;

\i /docker-entrypoint-initdb.d/01-schema.sql
\i /docker-entrypoint-initdb.d/03-triggers.sql
\i /docker-entrypoint-initdb.d/05-init-pg-ivm.sql


INSERT INTO manufacturers (name) VALUES ('BMW'), ('Audi'), ('Toyota'), ('Honda'), ('Aston Martin');
INSERT INTO models (name, manufacturer_id) VALUES ('M3', 1), ('A3', 2), ('Supra', 3), ('Civic', 4), ('DB9', 5);
SELECT 'Test database bd2_test created the schema successfully' AS status;
