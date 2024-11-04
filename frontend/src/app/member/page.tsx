'use client'
import CopyIcon from '@mui/icons-material/ContentCopy'
import DeleteIcon from '@mui/icons-material/Delete'
import { useUser } from '@/context/user'
import { client } from '@/utils/client'
import { Client } from '@/utils/types'
import {
  Box,
  Button,
  IconButton,
  List,
  ListItem,
  Typography,
} from '@mui/material'
import { blueGrey } from '@mui/material/colors'
import { useCallback, useEffect, useState } from 'react'
import Error from 'next/error'
import EditIcon from '@mui/icons-material/Edit'
import { useRouter } from 'next/navigation'

export default function Page() {
  const [clients, setClients] = useState<Client[]>([])
  const { user } = useUser()
  const fetchClients = useCallback(async () => {
    const { data, error } = await client.GET('/account/clients')
    if (error) {
      console.error(error)
      return
    }
    setClients(data.clients)
  }, [])
  const deleteClient = useCallback(
    async (clientId: string) => {
      const ok = confirm('Are you sure you want to delete this client?')
      if (!ok) {
        return
      }
      const { error } = await client.DELETE('/account/clients/{id}', {
        params: { path: { id: clientId } },
      })
      if (error) {
        console.error(error)
      }
      fetchClients()
    },
    [fetchClients]
  )
  const router = useRouter()
  const [copiedClient, setCopiedClient] = useState<string | null>(null)
  const [copiedTimer, setCopiedTimer] = useState<NodeJS.Timeout | null>(null)

  const copyClientId = useCallback(
    (clientId: string) => {
      navigator.clipboard.writeText(clientId)
      setCopiedClient(clientId)
      if (copiedTimer) {
        clearTimeout(copiedTimer)
      }
      const timer = setTimeout(() => {
        setCopiedClient(null)
      }, 3000)
      setCopiedTimer(timer)
    },
    [copiedTimer, setCopiedClient, setCopiedTimer]
  )

  useEffect(() => {
    fetchClients()
  }, [fetchClients])

  if (!user) {
    return <Error statusCode={500} />
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
        <List>
          {clients.map((client) => (
            <ListItem
              secondaryAction={
                <Box>
                  <IconButton
                    sx={{ marginX: 1 }}
                    onClick={() =>
                      router.push(`/member/clients/edit/${client.id}`)
                    }
                    edge="end"
                  >
                    <EditIcon />
                  </IconButton>
                  <IconButton
                    sx={{ marginX: 1 }}
                    onClick={() => deleteClient(client.id)}
                    edge="end"
                    aria-label="delete"
                  >
                    <DeleteIcon />
                  </IconButton>
                </Box>
              }
              key={client.id}
              sx={{
                bgcolor: 'white',
                paddingX: 4,
                paddingY: 2,
              }}
            >
              <Typography variant="h6">{client.name}</Typography>
              <Typography component="span" ml={2} variant="body1">
                ID:
                <Typography component="span">
                  <IconButton
                    onClick={() => copyClientId(client.id)}
                    size="small"
                    color={copiedClient === client.id ? 'success' : 'default'}
                  >
                    <CopyIcon />
                  </IconButton>
                </Typography>
                {client.id}
              </Typography>
            </ListItem>
          ))}
        </List>
      </Box>
    </Box>
  )
}
