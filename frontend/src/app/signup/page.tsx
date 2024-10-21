import { SignupForm } from '@/components/SignupForm'
import { Container, Typography } from '@mui/material'

export default function Page() {
  return (
    <Container component="main" sx={{ pt: 4, pb: 6 }}>
      <Typography
        variant="h4"
        fontWeight="bold"
        component="h1"
        gutterBottom
        mt={5}
      >
        ユーザー登録
      </Typography>
      <SignupForm />
    </Container>
  )
}
