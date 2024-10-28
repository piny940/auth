import createClient from 'openapi-fetch'
import { paths } from './api'

export const client = createClient<paths>({
  baseUrl: process.env.NEXT_PUBLIC_API_URL,
  credentials: 'include',
})
