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


const router = createBrowserRouter([
  {
    path: "/login", 
    element: <LoginPage/>
  },
  {
    element: <DashboardLayout/>,
    children: [
   {
    path: "/",
    element: <Dashboard/>
   }
    ]
  }
])

export const queryClient = new QueryClient({})
ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
  <Provider store={store}>
  <QueryClientProvider client={queryClient}>
   <RouterProvider router={router}/>
   </QueryClientProvider>
  </Provider>
  </React.StrictMode>,
)
