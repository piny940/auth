'use client'
import DeleteIcon from '@mui/icons-material/Delete'
import { useUser } from '@/context/user'
import { client } from '@/utils/client'
import { Client } from '@/utils/types'
import {
  Box,
  IconButton,
  List,
  ListItem,
  ListItemAvatar,
  Typography,
} from '@mui/material'
import { blueGrey } from '@mui/material/colors'
import { useCallback, useEffect, useState } from 'react'
import Error from 'next/error'

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
  useEffect(() => {
    fetchClients()
  }, [fetchClients])

  if (!user) {
    return <Error statusCode={500} />
  }
  return (
    <Box m={10} sx={{ '> *': { margin: 2 } }}>
      <Typography mb={4} variant="h4">
        ダッシュボード - {user.name}
      </Typography>
      <Box p={4} maxWidth="800px" borderRadius={2} bgcolor={blueGrey[50]}>
        <Typography mb={1} variant="h5">
          Clients
        </Typography>
        <List>
          {clients.map((client) => (
            <ListItem
              secondaryAction={
                <IconButton
                  onClick={() => deleteClient(client.id)}
                  edge="end"
                  aria-label="delete"
                >
                  <DeleteIcon />
                </IconButton>
              }
              key={client.id}
              sx={{ bgcolor: 'white' }}
            >
              <Typography variant="h6">{client.name}</Typography>
            </ListItem>
          ))}
        </List>
      </Box>
    </Box>
  )
}
