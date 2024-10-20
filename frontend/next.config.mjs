/** @type {import('next').NextConfig} */
import dotenv from 'dotenv'

if (process.env.NODE_ENV === 'development') {
  dotenv.config()
}

const nextConfig = {}

export default nextConfig
