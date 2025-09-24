import path from 'path'
import jwt from 'jsonwebtoken'
import { IDToken, User } from '@/utils/types'
import { verifyIDToken } from '@/utils/token'

type Props = {
  code: string
}
export default async function Page({
  searchParams,
}: {
  searchParams: Promise<Props>
}) {
  const query = await searchParams
  if (!process.env.NEXT_PUBLIC_API_URL) {
    throw new Error('NEXT_PUBLIC_API_URL is not set')
  }
  if (!process.env.NEXT_PUBLIC_CLIENT_ID) {
    throw new Error('NEXT_PUBLIC_CLIENT_ID is not set')
  }
  if (!process.env.NEXT_PUBLIC_APP_URL) {
    throw new Error('NEXT_PUBLIC_APP_URL is not set')
  }
  if (!process.env.CLIENT_SECRET) {
    throw new Error('CLIENT_SECRET is not set')
  }
  if (!process.env.OAUTH_RSA_PUBLIC_KEY) {
    throw new Error('OAUTH_RSA_PUBLIC_KEY is not set')
  }
  const secret = Buffer.from(
    `${process.env.NEXT_PUBLIC_CLIENT_ID}:${process.env.CLIENT_SECRET}`,
  )
  const res = await fetch(
    path.join(process.env.NEXT_PUBLIC_API_URL, 'oauth', 'token'),
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
        'Authorization': `Basic ${secret.toString('base64')}`,
      },
      body: new URLSearchParams({
        code: query.code,
        grant_type: 'authorization_code',
        redirect_uri: process.env.NEXT_PUBLIC_APP_URL + '/callback',
      }),
    },
  )
  if (!res.ok) {
    throw new Error(await res.text())
  }
  const json = await res.json()
  const idToken = jwt.verify(
    json.id_token,
    process.env.OAUTH_RSA_PUBLIC_KEY,
  ) as IDToken
  console.log(idToken)
  try {
    verifyIDToken(idToken)
  }
  catch (e) {
    throw new Error(e as string)
  }
  const user = idToken.sub
    .split(';')
    .map(kv => kv.split(':'))
    .reduce((acc, [k, v]) => ({ ...acc, [k]: v }), {}) as User
  return (
    <main className="m-20">
      <h1 className="text-5xl font-bold m-5">Successfully Authenticated!</h1>
      <p className="m-5">
        ID:
        {user.id}
      </p>
      <p className="m-5">
        Name:
        {user.name}
      </p>
    </main>
  )
}
