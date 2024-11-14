```mermaid
erDiagram
  users ||--o{ clients : "クライアント"
  clients ||--o{ redirect_uris : ""
  users ||--o{ auth_codes : "認可した"
  clients ||--o{ auth_codes : "認可された"
  auth_codes ||--o{ auth_code_scopes : ""
  users ||--o{ approvals : "承認した"
  clients ||--o{ approvals : "承認された"
  approvals ||--o{ approval_scopes : ""

  users {
    bigserial id PK
    varchar name
    varchar encrypted_password
    timestamptz created_at
    timestamptz updated_at
  }

  clients {
    uuid id PK
    varchar encrypted_secret
    bigint user_id FK
    varchar name
    timestamptz created_at
    timestamptz updated_at
  }
  redirect_uris {
    bigserial id PK
    varchar client_id FK
    varchar uri
    timestamptz created_at
    timestamptz updated_at
  }
  approvals {
    bigserial id PK
    bigint client_id FK
    bigint user_id FK
    timestamptz created_at
    timestamptz updated_at
  }
  approval_scopes {
    int scope_id PK
    bigint approval_id PK
    timestamptz created_at
    timestamptz updated_at
  }
  auth_codes {
    bigserial id PK
    varchar value
    varchar client_id FK
    bigint user_id FK
    varchar redirect_uri
    bool used
    timestamptz expires_at
    timestamptz auth_time
    timestamptz created_at
    timestamptz updated_at
  }
  auth_code_scopes {
    int scope_id PK
    bigint auth_code_id PK
    timestamptz created_at
    timestamptz updated_at
  }

```
