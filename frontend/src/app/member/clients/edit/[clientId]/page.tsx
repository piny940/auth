import { ClientEditPage } from '@/components/ClientEditPage'

type Props = {
  clientId: string
}
export default async function Page(req: { params: Promise<Props> }) {
  const params = await req.params
  return <ClientEditPage clientId={params.clientId} />
}
