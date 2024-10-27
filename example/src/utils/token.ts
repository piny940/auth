import { IDToken } from "./types"

export const verifyIDToken = (token: IDToken) => {
  if (token.iss !== process.env.NEXT_PUBLIC_ISSUER) {
    throw new Error("issuer is not valid")
  }
  if (token.aud !== process.env.NEXT_PUBLIC_CLIENT_ID) {
    throw new Error("audience is not valid")
  }
  if (token.exp < Date.now() / 1000) {
    throw new Error("token is expired")
  }
}
