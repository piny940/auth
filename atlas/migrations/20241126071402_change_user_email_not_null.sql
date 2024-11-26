-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "email" SET NOT NULL;
-- Create index "email" to table: "users"
CREATE UNIQUE INDEX "email" ON "public"."users" ("email");
