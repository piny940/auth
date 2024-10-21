'use client'

import { Controller, useForm } from 'react-hook-form'
import { Box, Button, TextField } from '@mui/material'
import { useCallback } from 'react'
import { client } from '@/utils/client'

type LoginInput = {
  name: string
  password: string
}
export const LoginForm = (): JSX.Element => {
  const { control, handleSubmit, setError } = useForm<LoginInput>({
    defaultValues: { name: '', password: '' },
  })
  const submit = useCallback(async (data: LoginInput) => {
    const { error } = await client.POST('/login', {
      body: { name: data.name, password: data.password },
    })
  }, [])
  return (
    <Box
      component="form"
      onSubmit={handleSubmit(submit)}
      sx={{ '> *': { margin: 2 } }}
    >
      <Box>
        <Controller
          name="name"
          control={control}
          render={({ field }) => (
            <TextField {...field} label="Name" variant="outlined" fullWidth />
          )}
        />
      </Box>
      <Box>
        <Controller
          name="password"
          control={control}
          render={({ field }) => (
            <TextField
              type="password"
              {...field}
              label="Password"
              variant="outlined"
              fullWidth
            />
          )}
        />
      </Box>
      <Box>
        <Button type="submit" fullWidth variant="contained">
          送信
        </Button>
      </Box>
    </Box>
  )
}
