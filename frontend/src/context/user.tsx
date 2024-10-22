'use client'
import { client } from '@/utils/client'
import { User } from '@/utils/types'
import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useState,
} from 'react'

interface UserContextInterface {
  user: User | null
  loading: boolean
  setUser: (user: User | null) => void
}

const defaultUserState: UserContextInterface = {
  user: null,
  loading: true,
  setUser: () => undefined,
}

const UserContext = createContext(defaultUserState)

const useUser = () => useContext(UserContext)

interface UserProviderProps {
  children: ReactNode
}

const UserProvider = ({ children }: UserProviderProps): JSX.Element => {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const value: UserContextInterface = { user, loading, setUser }

  useEffect(() => {
    client.GET('/me').then((res) => {
      setUser(res.data?.user || null)
      setLoading(false)
    })
  }, [])

  return <UserContext.Provider value={value}>{children}</UserContext.Provider>
}

export { useUser, UserProvider }
