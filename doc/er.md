```mermaid
erDiagram
  blogs ||--o| qiita_articles : "Qiita記事"

  clients {
    bigserial id PK
    bigint user_id
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
