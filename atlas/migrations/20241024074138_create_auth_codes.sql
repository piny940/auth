-- Modify "approval_scopes" table
ALTER TABLE "public"."approval_scopes" ADD COLUMN "created_at" timestamptz NOT NULL, ADD COLUMN "updated_at" timestamptz NOT NULL;
-- Create "auth_codes" table
CREATE TABLE "public"."auth_codes" ("id" bigserial NOT NULL, "value" character varying(32) NOT NULL, "client_id" character varying(16) NOT NULL, "user_id" bigint NOT NULL, "expires_at" timestamptz NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "client_id" FOREIGN KEY ("client_id") REFERENCES "public"."clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create "auth_code_scopes" table
CREATE TABLE "public"."auth_code_scopes" ("scope_id" integer NOT NULL, "auth_code_id" bigint NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("scope_id", "auth_code_id"), CONSTRAINT "auth_code_id" FOREIGN KEY ("auth_code_id") REFERENCES "public"."auth_codes" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
