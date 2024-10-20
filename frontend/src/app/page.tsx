import { LoginForm } from '@/components/LoginForm'
import { Container, Typography } from '@mui/material'

export default function Home() {
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
      <LoginForm />
    </Container>
  )
}
