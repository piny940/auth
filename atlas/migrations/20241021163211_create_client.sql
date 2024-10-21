-- Create "clients" table
CREATE TABLE "public"."clients" ("id" character varying(16) NOT NULL, "encrypted_secret" character varying(255) NOT NULL, "user_id" bigint NOT NULL, "name" character varying(255) NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "user_id" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "user_id" to table: "clients"
CREATE INDEX "user_id" ON "public"."clients" ("user_id");
