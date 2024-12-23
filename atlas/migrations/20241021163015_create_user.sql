-- Create "users" table
CREATE TABLE "public"."users" ("id" bigserial NOT NULL, "name" character varying(255) NOT NULL, "encrypted_password" character varying(255) NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create index "name" to table: "users"
CREATE UNIQUE INDEX "name" ON "public"."users" ("name");
