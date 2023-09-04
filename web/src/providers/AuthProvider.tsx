import React from 'react'
import { useQuery } from 'react-query'
import { getMe } from '../services/auth'
import { useDispatch } from 'react-redux'
import { init } from '../redux/authSlice'

function AuthProvider({children}: {children: React.ReactNode}) {
  const dispatch = useDispatch()
  useQuery(['user', 'me'], () => getMe(), {onSuccess(user) {
    dispatch(init({user: user}))
  },})
  return (
    <>
    {children}
    </>
  )
}

export default AuthProvider