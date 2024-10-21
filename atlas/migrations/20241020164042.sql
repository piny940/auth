-- Create "users" table
CREATE TABLE "public"."users" ("id" bigserial NOT NULL, "encrypted_password" character varying(255) NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL);
