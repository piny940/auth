# トークン仕様

## Access Token

```json
{
  "typ": "JWT",
  "alg": "RS256"
}
{
  "iss": "https://auth.piny940.com", // staging環境はhttps://stg-auth.piny940.com
  "exp": 0000000000,
  "iat": 0000000000,
  "jti": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // 32字のランダム文字列
  "scope": "openid profile email", // スペース区切り
}
```

## ID Token

```json
{
  "typ": "JWT",
  "alg": "RS256"
}
{
  "iss": "https://auth.piny940.com", // staging環境はhttps://stg-auth.piny940.com
  "sub": "{username}",
  "aud": "{client_id}",
  "exp": 0000000000,
  "iat": 0000000000, // トークン発行日時
}
```
