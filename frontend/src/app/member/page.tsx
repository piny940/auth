'use client'
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
import Link from 'next/link'

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
    <Box sx={{ '> *': { margin: 2 } }}>
      <Typography mb={4} variant="h4">
        ダッシュボード - {user.name}
      </Typography>
      <Box p={4} borderRadius={2} maxWidth={800} bgcolor={blueGrey[50]}>
        <Box display="flex" alignItems="center" mb={1}>
          <Typography mr={3} variant="h5">
            Clients
          </Typography>
          <Link href="/member/clients/new">
            <Box
              p={1}
              borderRadius={1}
              fontWeight={500}
              bgcolor={blueGrey[200]}
            >
              新規作成
            </Box>
          </Link>
        </Box>
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
