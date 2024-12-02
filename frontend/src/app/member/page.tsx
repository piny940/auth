'use client'
import { useUser } from '@/context/user'
import { Box, Button, Typography } from '@mui/material'
import { blueGrey } from '@mui/material/colors'
import { ClientList } from '@/components/ClientList'
import { ApprovalList } from '@/components/ApprovalList'

export default function Page() {
  const { user } = useUser()

  if (!user) {
    throw Promise.resolve()
  }

  return (
    <Box sx={{ '> *': { margin: 2 } }}>
      <Typography mb={4} variant="h4">
        ダッシュボード - {user.name}
      </Typography>
      <Box p={4} borderRadius={2} maxWidth={800} bgcolor={blueGrey[50]}>
        <Box display="flex" alignItems="center" mb={1}>
          <Typography mr={3} variant="h5">
            Clients
          </Typography>
          <Button href="/member/clients/new">新規作成</Button>
        </Box>
        <ClientList />
      </Box>
      <Box p={4} borderRadius={2} maxWidth={800} bgcolor={blueGrey[50]}>
        <Box display="flex" alignItems="center" mb={1}>
          <Typography mr={3} variant="h5">
            連携済みサービス
          </Typography>
        </Box>
        <ApprovalList />
      </Box>
    </Box>
  )
}
