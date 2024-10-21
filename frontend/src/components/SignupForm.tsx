'use client'
import { client } from '@/utils/client'
import { Box, Button, TextField, Typography } from '@mui/material'
import Link from 'next/link'
import { useCallback } from 'react'
import { Controller, useForm } from 'react-hook-form'

type SignupInput = {
  name: string
  password: string
  passwordConfirmation: string
}
export const SignupForm = (): JSX.Element => {
  const { control, handleSubmit, setError } = useForm<SignupInput>({
    defaultValues: { name: '', password: '' },
  })
  const submit = useCallback(
    async (data: SignupInput) => {
      const { error } = await client.POST('/signup', {
        body: {
          name: data.name,
          password: data.password,
          password_confirmation: data.passwordConfirmation,
        },
      })
      if (!!error) {
        setError('name', { message: error.error_description })
        setError('password', { message: error.error_description })
        return
      }
    },
    [setError]
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
              error={fieldState.invalid}
              helperText={fieldState.error?.message}
            />
          )}
        />
      </Box>
      <Box>
        <Controller
          name="passwordConfirmation"
          control={control}
          render={({ field, fieldState }) => (
            <TextField
              type="password"
              {...field}
              label="Password Confirmation"
              variant="outlined"
              fullWidth
              error={fieldState.invalid}
              helperText={fieldState.error?.message}
            />
          )}
        />
      </Box>
      <Link href="/">
        <Typography component="span" color="primary">
          すでにアカウントをお持ちの方はこちら
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
