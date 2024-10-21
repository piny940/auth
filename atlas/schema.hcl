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
  primary_key {
    columns = [column.id]
  }
  index "name" {
    columns = [column.name]
    unique = true
  }
}
table "clients" {
  schema = schema.public
  column "id" {
    type = varchar(16)
  }
  column "encrypted_secret" {
    type = varchar(255)
  }
  column "user_id" {
    type = bigint
  }
  column "name" {
    type = varchar(255)
  }
  column "created_at" {
    type = timestamptz
  }
  column "updated_at" {
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "user_id" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update = NO_ACTION
    on_delete = CASCADE
  }
  index "user_id" {
    columns = [column.user_id]
  }
}
table "redirect_uris" {
  schema = schema.public
  column "id" {
    type = bigserial
  }
  column "client_id" {
    type = varchar(16)
  }
  column "uri" {
    type = varchar(255)
  }
  column "created_at" {
    type = timestamptz
  }
  column "updated_at" {
    type = timestamptz
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "client_id" {
    columns = [column.client_id]
    ref_columns = [table.clients.column.id]
    on_update = NO_ACTION
    on_delete = CASCADE
  }
}
table "approvals" {
  schema = schema.public
  column "id" {
    type = bigserial
  }
  column "client_id" {
    type = varchar(16)
  }
  column "user_id" {
    type = bigint
  }
  column "scopes" {
    type = sql("int[]")
  }
  column "created_at" {
    type = timestamptz
  }
  column "updated_at" {
    type = timestamptz
  }
  foreign_key "client_id" {
    columns = [column.client_id]
    ref_columns = [table.clients.column.id]
    on_update = NO_ACTION
    on_delete = CASCADE
  }
  foreign_key "user_id" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update = NO_ACTION
    on_delete = CASCADE
  }
  index "user_id_client_id" {
    columns = [column.user_id, column.client_id]
    unique = true
  }
}
