-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS public;

-- Create "users" table
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  name text NULL
);
