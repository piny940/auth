-- Create "users" table
CREATE TABLE "public"."users" ("id" bigserial NOT NULL, "name" character varying(255) NOT NULL, "encrypted_password" character varying(255) NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create index "name" to table: "users"
CREATE UNIQUE INDEX "name" ON "public"."users" ("name");
-- Create "clients" table
CREATE TABLE "public"."clients" ("id" character varying(16) NOT NULL, "encrypted_secret" character varying(255) NOT NULL, "user_id" bigint NOT NULL, "redirect_uris" character varying(255)[] NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "user_id" to table: "clients"
CREATE INDEX "user_id" ON "public"."clients" ("user_id");
