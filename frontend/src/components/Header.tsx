'use client'
import { useUser } from '@/context/user'
import { client } from '@/utils/client'
import { AppBar, Box, Button, Toolbar, Typography } from '@mui/material'
import Link from 'next/link'
import { useCallback } from 'react'

export const Header = (): JSX.Element => {
  const { user } = useUser()
  const logout = useCallback(async () => {
    const { error } = await client.DELETE('/session')
    if (error) {
      console.error(error)
      return
    }
  }, [])
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
            <Button color="inherit" onClick={logout}>
              Logout
            </Button>
          )}
        </Box>
      </Toolbar>
    </AppBar>
  )
}
