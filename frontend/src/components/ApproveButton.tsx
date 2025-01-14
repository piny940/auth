'use client'
import { client } from '@/utils/client'
import { Button } from '@mui/material'
import { useRouter } from 'next/navigation'
import { useCallback, JSX } from 'react'

type ApproveButtonProps = {
  clientID: string
  scope: string
  next: string
}
export const ApproveButton = ({
  clientID,
  scope,
  next,
}: ApproveButtonProps): JSX.Element => {
  const router = useRouter()
  const approve = useCallback(async () => {
    const { error } = await client.POST('/account/approvals', {
      body: { client_id: clientID, scope: scope },
    })
    if (error) {
      throw new Error(error.error_description)
    }
    router.push(next)
  }, [clientID, scope, router, next])
  return (
    <Button onClick={approve} variant="contained">
      Approve
    </Button>
  )
}
