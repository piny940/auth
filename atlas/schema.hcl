schema "public"{}

table "users" {
  schema = schema.public
  column "id" {
    type = bigserial
  }
  column "email" {
    type = varchar(255)
  }
  column "email_verified" {
    type = boolean
    default = false
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
  index "email" {
    columns = [column.email]
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
table "approval_scopes" {
  schema = schema.public
  column "scope_id" {
    type = int
  }
  column "approval_id" {
    type = bigint
  }
  column "created_at" {
    type = timestamptz
  }
  column "updated_at" {
    type = timestamptz
  }
  primary_key {
    columns = [column.scope_id, column.approval_id]
  }
  foreign_key "approval_id" {
    columns = [column.approval_id]
    ref_columns = [table.approvals.column.id]
    on_update = NO_ACTION
    on_delete = CASCADE
  }
}
table "auth_codes" {
  schema = schema.public
  column "id" {
    type = bigserial
  }
  column "value" {
    type = varchar(32)
  }
  column "client_id" {
    type = varchar(16)
  }
  column "user_id" {
    type = bigint
  }
  column "redirect_uri" {
    type = varchar(255)
  }
  column "used" {
    type = boolean
  }
  column "expires_at" {
    type = timestamptz
  }
  column "auth_time" {
    type = timestamptz
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
  foreign_key "user_id" {
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update = NO_ACTION
    on_delete = CASCADE
  }
}
table "auth_code_scopes" {
  schema = schema.public
  column "scope_id" {
    type = int
  }
  column "auth_code_id" {
    type = bigint
  }
  column "created_at" {
    type = timestamptz
  }
  column "updated_at" {
    type = timestamptz
  }
  primary_key {
    columns = [column.scope_id, column.auth_code_id]
  }
  foreign_key "auth_code_id" {
    columns = [column.auth_code_id]
    ref_columns = [table.auth_codes.column.id]
    on_update = NO_ACTION
    on_delete = CASCADE
  }
}
