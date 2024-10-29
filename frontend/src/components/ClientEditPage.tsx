'use client'
import { client } from '@/utils/client'
import { Client } from '@/utils/types'
import { useEffect, useState } from 'react'
import { ClientEdit } from './ClientEdit'

type ClientEditProps = {
  clientId: string
}
export const ClientEditPage = ({ clientId }: ClientEditProps) => {
  const [current, setCurrent] = useState<Client | null>(null)
  const fetchClient = async (clientId: string) => {
    const { data, error } = await client.GET('/account/clients/{id}', {
      params: { path: { id: clientId } },
    })
    if (error) {
      console.error(error)
      throw new Error('Failed to fetch client: ' + error)
    }
    setCurrent(data.client)
  }

  useEffect(() => {
    fetchClient(clientId)
  }, [clientId])

  if (!current) {
    return <div>Loading...</div>
  }
  return <ClientEdit client={current} />
}
