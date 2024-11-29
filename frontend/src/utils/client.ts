import createClient from 'openapi-fetch'
import { paths } from './api'
import { fromCookie } from './cookie'

export const client = createClient<paths>({
  fetch: (input: Request) => {
    const csrf = fromCookie('_csrf')
    if (csrf) {
      input.headers.set('X-CSRF-Token', csrf)
    }
    return fetch(input)
  },
  baseUrl: process.env.NEXT_PUBLIC_API_URL,
  credentials: 'include',
})
