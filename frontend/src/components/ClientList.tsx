import EditIcon from '@mui/icons-material/Edit'
import DeleteIcon from '@mui/icons-material/Delete'
import CopyIcon from '@mui/icons-material/FileCopy'
import { Client } from '@/utils/types'
import { Box, IconButton, List, ListItem, Typography } from '@mui/material'
import { useRouter } from 'next/navigation'
import { useCallback, useEffect, useState } from 'react'
import { useUser } from '@/context/user'
import { client } from '@/utils/client'

type ClientListProps = {}
export const ClientList = ({}: ClientListProps): JSX.Element => {
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

  if (!user || !clients) {
    throw Promise.resolve()
  }

  return (
    <List>
      {clients.map((client) => (
        <ListItem
          secondaryAction={
            <Box>
              <IconButton
                sx={{ marginX: 1 }}
                onClick={() => router.push(`/member/clients/edit/${client.id}`)}
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
  )
}
