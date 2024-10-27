import Link from "next/link"

export default function Home() {
  if (!process.env.NEXT_PUBLIC_APP_URL) {
    throw new Error("NEXT_PUBLIC_APP_URL is not set")
  }
  if (!process.env.NEXT_PUBLIC_API_URL) {
    throw new Error("NEXT_PUBLIC_API_URL is not set")
  }
  const redirectUri = process.env.NEXT_PUBLIC_APP_URL + "/callback"
  const link = `${process.env.NEXT_PUBLIC_API_URL}/oauth/authorize?client_id=${process.env.NEXT_PUBLIC_CLIENT_ID}&response_type=code&redirect_uri=${redirectUri}&scope=openid`
  return (
    <main className="m-20">
      <h1 className="text-4xl font-bold">Example Client</h1>
      <Link
        href={link}
        className="inline-block m-10 bg-sky-400 rounded p-3 text-lg font-medium hover:bg-sky-500"
      >
        mikanでログイン
      </Link>
    </main>
  )
}
