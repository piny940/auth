# Auth

OAuth2.0、OpenID Connect

## 開発

dependency: go, task, aqua

- task install
- task dev

### DB

マイグレーションファイル作成：`cd atlas && atlas migrate diff {migration name} --env local`

マイグレーション実行：`cd atlas && atlas migrate apply --env local`

DB リセット：`cd atlas && atlas schema clean --env local`
