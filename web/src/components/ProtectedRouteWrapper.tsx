import  { useEffect } from 'react'
import { Outlet, useNavigate } from 'react-router-dom'
import { useAppSelector } from '../redux/store'
import { ArrowClockwise } from 'react-bootstrap-icons'

function ProtectedRouteWrapper() {
    const {isLoggedIn, init} = useAppSelector(state => state.auth)
    const navigate = useNavigate()

    useEffect(() => {
        if(isLoggedIn && init){
            navigate("/")
        }

        if(!isLoggedIn && init){
            navigate("/login")
        }

    }, [isLoggedIn, init, navigate])

    if (!init || !isLoggedIn){
        return <div className='w-full py-4 flex justify-center items-center'>
            <ArrowClockwise className='animate-spin h-5 w-5'/>
        </div>
    }


  return (
    <>
    <Outlet/>
    </>
  )
}

export default ProtectedRouteWrapper