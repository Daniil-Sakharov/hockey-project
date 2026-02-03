import { QueryProvider, RouterProvider } from './providers'
import './styles/globals.css'

export function App() {
  return (
    <QueryProvider>
      <RouterProvider />
    </QueryProvider>
  )
}
