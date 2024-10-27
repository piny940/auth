# Auth

OAuth2.0・OpenID Connect プロバイダを作ってみた

ER 図: [docs/er.md](docs/er.md)
API スキーマ：<http://piny940.github.io/auth/>

## 開発

- dependency: go, task, aqua, db 作成
- .env ファイル作成
- .env.test ファイル作成

- task install
- task dev

### DB

- マイグレーションファイル作成：`task migrate:create -- {name}`
  - 作成取り消し：`task migrate:remove`
- マイグレーション実行：`task migrate`
  - マイグレーション取り消し：`task migrate:down`
