import { client } from '@/utils/client'
import { Container, Typography } from '@mui/material'

type Props = {
  searchParams: {
    next?: string
    scope: string
    client_id: string
  }
}
export default async function Page({ searchParams: query }: Props) {
  const { data, error } = await client.GET('/clients/{id}', {
    params: { path: { id: query.client_id } },
  })
  if (error) {
    throw new Error(error.error_description)
  }
  return (
    <Container component="main" sx={{ pt: 4, pb: 6 }}>
      <Typography
        variant="h4"
        fontWeight="bold"
        component="h1"
        gutterBottom
        mt={5}
      >
        {data.client.name} にアクセスを許可しますか？
      </Typography>
    </Container>
  )
}
