import { Approval } from '@/utils/types'
import { Button, List, ListItem, Typography } from '@mui/material'
import { useRouter } from 'next/navigation'
import { useCallback, useEffect, useState, JSX } from 'react'
import { useUser } from '@/context/user'
import { client } from '@/utils/client'

type ApprovalListProps = {}
export const ApprovalList = ({}: ApprovalListProps): JSX.Element => {
  const [approvals, setApprovals] = useState<Approval[]>([])
  const { user } = useUser()
  const fetchApprovals = useCallback(async () => {
    const { data, error } = await client.GET('/account/approvals')
    if (error) {
      console.error(error)
      return
    }
    setApprovals(data.approvals)
  }, [])
  const deleteApproval = useCallback(
    async (approvalId: number) => {
      const ok = confirm('Are you sure you want to delete this approval?')
      if (!ok) {
        return
      }
      const { error } = await client.DELETE('/account/approvals/{id}', {
        params: { path: { id: approvalId } },
      })
      if (error) {
        console.error(error)
      }
      fetchApprovals()
    },
    [fetchApprovals]
  )
  const router = useRouter()

  useEffect(() => {
    fetchApprovals()
  }, [fetchApprovals])

  if (!user || !approvals) {
    throw Promise.resolve()
  }

  return (
    <List>
      {approvals.map((approval) => (
        <ListItem
          secondaryAction={
            <Button
              onClick={() => deleteApproval(approval.id)}
              size="small"
              variant="outlined"
              color="error"
            >
              連携解除
            </Button>
          }
          key={approval.id}
          sx={{
            bgcolor: 'white',
            paddingX: 4,
            paddingY: 2,
          }}
        >
          <Typography variant="h6">{approval.client.name}</Typography>
          <Typography component="span" mx={2}>
            権限:
            <Typography mx={1} component="span">
              {approval.scopes.join(', ')}
            </Typography>
          </Typography>
        </ListItem>
      ))}
    </List>
  )
}
