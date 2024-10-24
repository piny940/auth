import { SignupForm } from '@/components/SignupForm'
import { Container, Typography } from '@mui/material'

type Props = {
  searchParams: Promise<{
    next?: string
  }>
}

export default async function Page(props: Props) {
  const query = await props.searchParams
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
      <SignupForm next={query.next || '/'} />
    </Container>
  )
}
