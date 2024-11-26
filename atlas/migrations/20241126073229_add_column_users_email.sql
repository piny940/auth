-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "email" character varying(255) NULL, ADD COLUMN "email_verified" boolean NULL DEFAULT false;
