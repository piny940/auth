# Auth

OAuth2.0・OpenID Connect プロバイダを作ってみた

- ER 図: [docs/er.md](docs/er.md)
- API スキーマ：<https://auth-doc.piny940.com/>

## 開発

- dependency: go, task, aqua, db 作成
- .env ファイル作成
- .env.test ファイル作成

- task install
- task dev

### DB

- マイグレーション実行：`task migrate:up`
- マイグレーション取り消し：`task migrate:down`
