'use client'
import { useUser } from '@/context/user'
import { client } from '@/utils/client'
import { Box, Button, TextField, Typography } from '@mui/material'
import Link from 'next/link'
import { useRouter, useSearchParams } from 'next/navigation'
import { useCallback } from 'react'
import { Controller, useForm } from 'react-hook-form'

type SignupInput = {
  name: string
  email: string
  password: string
  passwordConfirmation: string
}
export type SignupFormProps = {
  next: string
}
export const SignupForm = ({ next }: SignupFormProps): JSX.Element => {
  const { refresh } = useUser()
  const router = useRouter()
  const query = useSearchParams()
  const { control, handleSubmit, setError } = useForm<SignupInput>({
    defaultValues: {
      name: '',
      email: '',
      password: '',
      passwordConfirmation: '',
    },
  })
  const submit = useCallback(
    async (data: SignupInput) => {
      const { error } = await client.POST('/users/signup', {
        body: {
          name: data.name,
          email: data.email,
          password: data.password,
          password_confirmation: data.passwordConfirmation,
        },
      })
      if (!!error) {
        if (error.error === 'name_already_used') {
          setError('name', { message: error.error })
        } else if (
          error.error === 'email_already_used' ||
          error.error === 'email_format_invalid'
        ) {
          setError('email', { message: error.error_description })
        } else if (error.error === 'password_length_not_enough') {
          setError('password', { message: error.error_description })
        } else if (error.error === 'password_confirmation_not_match') {
          setError('password', { message: error.error_description })
          setError('passwordConfirmation', { message: error.error_description })
        } else if (error.error === 'name_length_not_enough') {
          setError('name', { message: error.error_description })
        } else {
          setError('name', { message: error.error_description })
          setError('password', { message: error.error_description })
          setError('passwordConfirmation', { message: error.error_description })
        }
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
          name="email"
          control={control}
          render={({ field, fieldState }) => (
            <TextField
              label="Email"
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
              required
              error={fieldState.invalid}
              helperText={fieldState.error?.message}
            />
          )}
        />
      </Box>
      <Link href={`/?${query.toString()}`}>
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
