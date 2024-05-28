-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS public;

-- Create "users" table
CREATE TABLE users (
  id integer NOT NULL,
  name text NULL,
  PRIMARY KEY (id)
);
