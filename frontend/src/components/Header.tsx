'use client'
import { useUser } from '@/context/user'
import { client } from '@/utils/client'
import { AppBar, Box, Button, Toolbar, Typography } from '@mui/material'
import Link from 'next/link'
import { useCallback, JSX } from 'react'

export const Header = (): JSX.Element => {
  const { user, refresh } = useUser()
  const logout = useCallback(async () => {
    const { error } = await client.DELETE('/session')
    if (error) {
      console.error(error)
      return
    }
    refresh()
  }, [refresh])
  return (
    <AppBar position="sticky">
      <Toolbar sx={{ display: 'flex', justifyContent: 'space-between' }}>
        <Link href="/" className="unstyled">
          <Typography variant="h5" fontWeight="bold">
            mikan
          </Typography>
        </Link>
        <Box pr={4}>
          {user && (
            <Box>
              <Button sx={{ mx: 1 }} color="inherit" href="/member">
                Dashboard
              </Button>
              <Button sx={{ mx: 1 }} color="inherit" onClick={logout}>
                Logout
              </Button>
            </Box>
          )}
        </Box>
      </Toolbar>
    </AppBar>
  )
}
