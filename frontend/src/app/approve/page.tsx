import { ApproveButton } from '@/components/ApproveButton'
import { client } from '@/utils/client'
import { Box, Container, Typography } from '@mui/material'

type Props = {
  searchParams: Promise<{
    next: string
    scope: string
    client_id: string
  }>
}
export default async function Page(props: Props) {
  const query = await props.searchParams
  if (!query.scope || !query.client_id || !query.next) {
    throw new Error('invalid query')
  }
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
        Will you authorize client:
        {' '}
        {data.client.name}
        {' '}
        ?
      </Typography>
      <Box sx={{ '> *': { margin: 2 } }}>
        <Typography component="p" mt={2}>
          Scope:
          {' '}
          {query.scope}
        </Typography>
        <ApproveButton
          clientID={query.client_id}
          scope={query.scope}
          next={query.next}
        />
      </Box>
    </Container>
  )
}
