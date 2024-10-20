'use client'

import { Controller, useForm } from 'react-hook-form'
import { Box, Button, TextField } from '@mui/material'

type LoginInput = {
  id: string
  password: string
}
export const LoginForm = (): JSX.Element => {
  const { control, getValues, handleSubmit, setError } = useForm<LoginInput>({
    defaultValues: { id: '', password: '' },
  })
  return (
    <Box component="form" sx={{ '> *': { margin: 2 } }}>
      <Box>
        <Controller
          name="id"
          control={control}
          render={({ field }) => (
            <TextField {...field} label="ID" variant="outlined" fullWidth />
          )}
        />
      </Box>
      <Box>
        <Controller
          name="password"
          control={control}
          render={({ field }) => (
            <TextField
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
