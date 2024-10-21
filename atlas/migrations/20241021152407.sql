-- Create "users" table
CREATE TABLE "public"."users" ("id" bigserial NOT NULL, "name" character varying(255) NOT NULL, "encrypted_password" character varying(255) NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create index "name" to table: "users"
CREATE UNIQUE INDEX "name" ON "public"."users" ("name");
-- Create "clients" table
CREATE TABLE "public"."clients" ("id" character varying(16) NOT NULL, "encrypted_secret" character varying(255) NOT NULL, "user_id" bigint NOT NULL, "name" character varying(255) NOT NULL, "redirect_uris" character varying(255)[] NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "user_id" to table: "clients"
CREATE INDEX "user_id" ON "public"."clients" ("user_id");
-- Create "approvals" table
CREATE TABLE "public"."approvals" ("id" bigserial NOT NULL, "client_id" character varying(16) NOT NULL, "user_id" bigint NOT NULL, "scopes" integer[] NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, CONSTRAINT "client_id" FOREIGN KEY ("client_id") REFERENCES "public"."clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "user_id_client_id" to table: "approvals"
CREATE UNIQUE INDEX "user_id_client_id" ON "public"."approvals" ("user_id", "client_id");
