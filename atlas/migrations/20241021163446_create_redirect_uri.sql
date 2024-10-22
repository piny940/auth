-- Create "redirect_uris" table
CREATE TABLE "public"."redirect_uris" ("id" bigserial NOT NULL, "client_id" character varying(16) NOT NULL, "uri" character varying(255) NOT NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "client_id" FOREIGN KEY ("client_id") REFERENCES "public"."clients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
