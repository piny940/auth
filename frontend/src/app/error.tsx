'use client'

import { Container, Typography } from '@mui/material'
import Link from 'next/link'

export default function GlobalError({
  error,
}: {
  error: Error & { digest?: string }
  reset: () => void
}) {
  return (
    <Container component="main" sx={{ pt: 4, pb: 6 }}>
      <Typography
        variant="h4"
        fontWeight="bold"
        component="h1"
        gutterBottom
        mt={5}
      >
        {error.message}
      </Typography>
      <Link href="/">
        <Typography variant="h6" color="primary">
          Back to Home
        </Typography>
      </Link>
    </Container>
  )
}
