import createClient from 'openapi-fetch'
import { paths } from './api'

export const client = createClient<paths>({
  baseUrl: process.env.BACKEND_URL + '/api/v1',
})
