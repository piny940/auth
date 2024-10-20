'use client'

import { Controller, useForm } from 'react-hook-form'
import styles from './page.module.css'
import { Box, Button, Container, TextField, Typography } from '@mui/material'

type LoginInput = {
  id: string
  password: string
}

export default function Home() {
  const { control, getValues, handleSubmit, setError } = useForm<LoginInput>({
    defaultValues: { id: '', password: '' },
  })
  return (
    <Container component="main" sx={{ pt: 4, pb: 6 }}>
      <Typography
        variant="h4"
        fontWeight="bold"
        component="h1"
        gutterBottom
        mt={5}
      >
        ログイン
      </Typography>
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
    </Container>
  )
}
