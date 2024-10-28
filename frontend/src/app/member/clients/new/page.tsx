'use client'
import {
  RedirectURIsEdit,
  RedirectURIsFields,
} from '@/components/RedirectURIsEdit'
import { client } from '@/utils/client'
import { Box, Button, TextField, Typography } from '@mui/material'
import { useCallback } from 'react'
import { Controller, useForm } from 'react-hook-form'

type ClientInput = {
  name: string
  redirectUris: string[]
}

export default function Page() {
  const { handleSubmit, getValues, control } = useForm<ClientInput>({
    defaultValues: {
      name: '',
      redirectUris: [],
    },
  })
  const { control: redirectURIsControl, getValues: getRedirectURIsValues } =
    useForm<RedirectURIsFields>({
      defaultValues: {
        redirectURIs: [],
      },
    })
  const requiredRule = { required: 'このフィールドは必須です。' }
  const submit = useCallback(async () => {
    const c = getValues()
    const redirectURIs = getRedirectURIsValues()
    const { error } = await client.POST('/account/clients', {
      body: {
        client: {
          name: c.name,
          redirect_urls: redirectURIs.redirectURIs.map((r) => r.url),
        },
      },
    })
    if (error) {
      console.error(error)
      return
    }
  }, [getValues, getRedirectURIsValues])
  return (
    <Box>
      <Typography variant="h4">New Client</Typography>
      <Box
        onSubmit={handleSubmit(submit)}
        component="form"
        sx={{ '> *': { margin: 2 } }}
      >
        <Box>
          <Controller
            control={control}
            name="name"
            rules={requiredRule}
            render={({ field, fieldState }) => (
              <TextField
                fullWidth
                label="Name"
                error={fieldState.invalid}
                helperText={fieldState.error?.message}
                {...field}
              />
            )}
          />
        </Box>
        <RedirectURIsEdit control={redirectURIsControl} />
        <Box>
          <Button type="submit" fullWidth variant="contained">
            Submit
          </Button>
        </Box>
      </Box>
    </Box>
  )
}
