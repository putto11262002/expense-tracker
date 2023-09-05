import React from 'react'
import ReactDOM from 'react-dom/client'

import './index.css'
import { QueryClient, QueryClientProvider } from 'react-query'
import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import LoginPage from './pages/LoginPage.tsx'
import DashboardLayout from './components/layout/DashboardLayout.tsx'
import Dashboard from './pages/Dashboard.tsx'
import { Provider } from 'react-redux'
import store from './redux/store.ts'
import InitAuth from './components/InitAuth.tsx'
import ProtectedRouteWrapper from './components/ProtectedRouteWrapper.tsx'
import RegisterPage from './pages/RegisterPage.tsx'


const router = createBrowserRouter([
  {
    path: "/login", 
    element: <LoginPage/>
  },
  {
    path: "/register",
    element: <RegisterPage/>
  },
  {
    element: <DashboardLayout/>,
  children: [
    {
      element: <ProtectedRouteWrapper/>,
      children: [
        {
         path: "/",
         element: <Dashboard/>
        }
         ]
    }
  ]
  }
])

const queryClient = new QueryClient({})
ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
  <Provider store={store}>
  <QueryClientProvider client={queryClient}>
   <InitAuth>
   <RouterProvider router={router}/>
   </InitAuth>
   </QueryClientProvider>
  </Provider>
  </React.StrictMode>,
)
