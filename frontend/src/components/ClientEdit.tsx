'use client'

import { useForm } from 'react-hook-form'
import { ClientForm, ClientInput } from './ClientForm'
import { RedirectURIsFields } from './RedirectURIsEdit'
import { Box, Typography } from '@mui/material'
import { Client } from '@/utils/types'

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
  const { control: redirectURIsControl, getValues: getRedirectURIsValues } =
    useForm<RedirectURIsFields>({
      defaultValues: {
        redirectURIs: client.redirect_urls.map((url) => ({ url })),
      },
    })

  return (
    <Box>
      <Typography variant="h4">Edit Client</Typography>
      <ClientForm
        control={control}
        redirectURIsControl={redirectURIsControl}
        submit={handleSubmit(() => {
          console.log(getValues())
          console.log(getRedirectURIsValues())
        })}
      />
    </Box>
  )
}
