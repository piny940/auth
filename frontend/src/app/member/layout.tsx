'use client'
import { useUser } from '@/context/user'
import { Box } from '@mui/material'
import { useRouter } from 'next/navigation'
import { ReactNode } from 'react'

export default function Layout({ children }: { children: ReactNode }) {
  const { user, loading } = useUser()
  const router = useRouter()

  if (loading) {
    return <Box>Loading...</Box>
  }
  if (!user) {
    router.push('/')
  }
  return <Box>{children}</Box>
}
