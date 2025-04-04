'use client'
import CopyIcon from '@mui/icons-material/ContentCopy'
import { ClientForm, ClientInput } from '@/components/ClientForm'
import { RedirectURIsFields } from '@/components/RedirectURIsEdit'
import { client } from '@/utils/client'
import { Box, IconButton, Typography } from '@mui/material'
import { blue } from '@mui/material/colors'
import Link from 'next/link'
import { useCallback, useState } from 'react'
import { useForm } from 'react-hook-form'

export default function Page() {
  const { handleSubmit, getValues, control } = useForm<ClientInput>({
    defaultValues: {
      name: '',
      redirectUris: [],
    },
  })
  const { control: redirectURIsControl, getValues: getRedirectURIsValues } =
    useForm<RedirectURIsFields>({
      defaultValues: {
        redirectURIs: [],
      },
    })
  const [created, setCreated] = useState<{ id: string; secret: string } | null>(
    null
  )
  const submit = useCallback(async () => {
    const c = getValues()
    const redirectURIs = getRedirectURIsValues()
    const { data, error } = await client.POST('/account/clients', {
      body: {
        client: {
          name: c.name,
          redirect_urls: redirectURIs.redirectURIs.map((r) => r.url),
        },
      },
    })
    if (error) {
      console.error(error)
      return
    }
    setCreated({ id: data.client.id, secret: data.client.secret })
  }, [getValues, getRedirectURIsValues])
  return (
    <Box>
      <Typography variant="h4">New Client</Typography>
      {created ? (
        <Box m={2}>
          <Typography variant="h5">Clientを作成しました。</Typography>
          <Typography component="p">
            Secretは後で確認できません。必ず保存してください。
          </Typography>
          <Box m={2}>
            <Typography>
              Client ID:
              <Typography component="span">
                <IconButton
                  onClick={() => navigator.clipboard.writeText(created.id)}
                >
                  <CopyIcon />
                </IconButton>
                {created.id}
              </Typography>
            </Typography>
            <Typography>
              Client Secret:
              <Typography component="span">
                <IconButton
                  onClick={() => navigator.clipboard.writeText(created.secret)}
                >
                  <CopyIcon />
                </IconButton>
                {created.secret}
              </Typography>
            </Typography>
          </Box>
          <Link href="/member">
            <Typography sx={{ color: blue[700] }}>
              Back to member page
            </Typography>
          </Link>
        </Box>
      ) : (
        <ClientForm
          control={control}
          redirectURIsControl={redirectURIsControl}
          submit={handleSubmit(submit)}
        />
      )}
    </Box>
  )
}
