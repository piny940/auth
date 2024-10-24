-- Create "approvals" table
CREATE TABLE "public"."approvals" ("id" bigserial NOT NULL, "client_id" character varying(16) NOT NULL, "user_id" bigint NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "client_id" FOREIGN KEY ("client_id") REFERENCES "public"."clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "user_id_client_id" to table: "approvals"
CREATE UNIQUE INDEX "user_id_client_id" ON "public"."approvals" ("user_id", "client_id");
-- Create "approval_scopes" table
CREATE TABLE "public"."approval_scopes" ("scope_id" integer NOT NULL, "approval_id" bigint NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("scope_id", "approval_id"), CONSTRAINT "approval_id" FOREIGN KEY ("approval_id") REFERENCES "public"."approvals" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
