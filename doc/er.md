```mermaid
erDiagram
  blogs ||--o| qiita_articles : "Qiita記事"

  users {
    bigserial id PK
    varchar name
    varchar password
  }

  clients {
    bigserial id PK
    bigint user_id FK
    varchar name
    varchar[] redirect_urls
  }
  tokens {
    bigserial id PK
    bigint client_id
    varchar value
    varchar refresh_token
    timestamp expires_at
    varchar auth_code
    int[] scopes
  }

```
