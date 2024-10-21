schema "public"{}

table "users" {
  schema = schema.public
  column "id" {
    type = bigserial
  }
  column "name" {
    type = varchar(255)
  }
  column "encrypted_password" {
    type = varchar(255)
  }
  column "created_at" {
    type = timestamptz
  }
  column "updated_at" {
    type = timestamptz
  }
  index "users_name" {
    columns = [column.name]
    unique = true
  }
}
