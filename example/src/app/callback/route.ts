import { NextRequest } from "next/server"
import path from "path"
import jwt from "jsonwebtoken"
import { IDToken } from "@/utils/types"
import { verifyIDToken } from "@/utils/token"

export async function GET(request: NextRequest) {
  if (!process.env.NEXT_PUBLIC_API_URL) {
    throw new Error("NEXT_PUBLIC_API_URL is not set")
  }
  const code = request.nextUrl.searchParams.get("code")
  if (!code) {
    throw new Error("code is not set")
  }
  if (!process.env.NEXT_PUBLIC_CLIENT_ID) {
    throw new Error("NEXT_PUBLIC_CLIENT_ID is not set")
  }
  if (!process.env.NEXT_PUBLIC_APP_URL) {
    throw new Error("NEXT_PUBLIC_APP_URL is not set")
  }
  if (!process.env.CLIENT_SECRET) {
    throw new Error("CLIENT_SECRET is not set")
  }
  if (!process.env.SERVER_PUBLIC_KEY) {
    throw new Error("SERVER_PUBLIC_KEY is not set")
  }
  const secret = Buffer.from(
    `${process.env.NEXT_PUBLIC_CLIENT_ID}:${process.env.CLIENT_SECRET}`
  )
  const res = await fetch(
    path.join(process.env.NEXT_PUBLIC_API_URL, "oauth", "token"),
    {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
        Authorization: `Basic ${secret.toString("base64")}`,
      },
      body: new URLSearchParams({
        client_id: process.env.NEXT_PUBLIC_CLIENT_ID,
        code: code,
        grant_type: "authorization_code",
        redirect_uri: process.env.NEXT_PUBLIC_APP_URL + "/callback",
      }),
    }
  )
  const json = await res.json()
  const idToken = jwt.verify(
    json.id_token,
    process.env.SERVER_PUBLIC_KEY
  ) as IDToken
  try {
    verifyIDToken(idToken)
  } catch (e) {
    throw new Error(e as string)
  }
  return new Response()
}
