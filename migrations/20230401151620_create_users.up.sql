CREATE TABLE IF NOT EXISTS "users" (
  id serial NOT NULL PRIMARY KEY,
  name character varying(300) NOT NULL unique,
  balance DECIMAL(10,2) NOT NULL
);