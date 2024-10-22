'use client'

import { Controller, useForm } from 'react-hook-form'
import { Box, Button, TextField, Typography } from '@mui/material'
import { useCallback } from 'react'
import { client } from '@/utils/client'
import Link from 'next/link'
import { useUser } from '@/context/user'
import { useRouter, useSearchParams } from 'next/navigation'

type LoginInput = {
  name: string
  password: string
}
type LoginFormProps = {
  next: string
}
export const LoginForm = ({ next }: LoginFormProps): JSX.Element => {
  const { refresh } = useUser()
  const router = useRouter()
  const query = useSearchParams()
  const { control, handleSubmit, setError } = useForm<LoginInput>({
    defaultValues: { name: '', password: '' },
  })
  const submit = useCallback(
    async (data: LoginInput) => {
      const { error } = await client.POST('/session', {
        body: { name: data.name, password: data.password },
      })
      if (!!error) {
        setError('name', { message: error.error_description })
        setError('password', { message: error.error_description })
        return
      }
      refresh()
      router.push(next)
    },
    [setError, refresh, router, next]
  )

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
          render={({ field, fieldState }) => (
            <TextField
              label="Name"
              variant="outlined"
              fullWidth
              required
              error={fieldState.invalid}
              helperText={fieldState.error?.message}
              {...field}
            />
          )}
        />
      </Box>
      <Box>
        <Controller
          name="password"
          control={control}
          render={({ field, fieldState }) => (
            <TextField
              type="password"
              {...field}
              label="Password"
              variant="outlined"
              fullWidth
              required
              error={fieldState.invalid}
              helperText={fieldState.error?.message}
            />
          )}
        />
      </Box>
      <Link href={`/signup?${query.toString()}`}>
        <Typography component="span" color="primary">
          新規アカウント登録
        </Typography>
      </Link>
      <Box>
        <Button type="submit" fullWidth variant="contained">
          送信
        </Button>
      </Box>
    </Box>
  )
}
