'use client'
import { client } from '@/utils/client'
import { User } from '@/utils/types'
import {
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useEffect,
  useState,
  JSX,
} from 'react'

interface UserContextInterface {
  user: User | null
  loading: boolean
  setUser: (user: User | null) => void
  refresh: () => void
}

const defaultUserState: UserContextInterface = {
  user: null,
  loading: true,
  setUser: () => undefined,
  refresh: () => undefined,
}

const UserContext = createContext(defaultUserState)

const useUser = () => useContext(UserContext)

interface UserProviderProps {
  children: ReactNode
}

const UserProvider = ({ children }: UserProviderProps): JSX.Element => {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const refresh = useCallback(async () => {
    setLoading(true)
    const res = await client.GET('/session')
    setUser(res.data?.user || null)
    setLoading(false)
  }, [setLoading, setUser])

  const value: UserContextInterface = { user, loading, setUser, refresh }

  useEffect(() => {
    refresh()
  }, [refresh])

  return <UserContext.Provider value={value}>{children}</UserContext.Provider>
}

export { useUser, UserProvider }
