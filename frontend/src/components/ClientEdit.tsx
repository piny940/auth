'use client'

import { useForm } from 'react-hook-form'
import { ClientForm, ClientInput } from './ClientForm'
import { RedirectURIsFields } from './RedirectURIsEdit'
import { Box, Typography } from '@mui/material'
import { Client } from '@/utils/types'
import { useCallback } from 'react'
import { client as apiClient } from '@/utils/client'
import { useRouter } from 'next/navigation'

export type ClientEditProps = {
  client: Client
}

export const ClientEdit = ({ client }: ClientEditProps) => {
  const { handleSubmit, getValues, control } = useForm<ClientInput>({
    defaultValues: {
      name: client.name,
      redirectUris: client.redirect_urls,
    },
  })
  const { control: redirectURIsControl, getValues: getRedirectURIsValues }
    = useForm<RedirectURIsFields>({
      defaultValues: {
        redirectURIs: client.redirect_urls.map(url => ({ url })),
      },
    })
  const router = useRouter()
  const submit = useCallback(async () => {
    const newClient = getValues()
    const urls = getRedirectURIsValues().redirectURIs.map(uri => uri.url)
    const { error } = await apiClient.POST('/account/clients/{id}', {
      params: { path: { id: client.id } },
      body: {
        client: {
          name: newClient.name,
          redirect_urls: urls,
        },
      },
    })
    if (error) {
      console.error(error)
      throw new Error('Failed to update client: ' + error)
    }
    router.push('/member')
  }, [getValues, getRedirectURIsValues, client.id, router])

  return (
    <Box>
      <Typography variant="h4">Edit Client</Typography>
      <ClientForm
        control={control}
        redirectURIsControl={redirectURIsControl}
        submit={handleSubmit(submit)}
      />
    </Box>
  )
}
