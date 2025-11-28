import { StrictMode, useEffect, useState } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { api } from './services/api.ts'
import { LoginPage } from './components/LoginPage.tsx'



const Root = () => {
  const [isAuthenticated, setAuthenticated] = useState<boolean|null>(null)

  useEffect(() => {
    api.me()
      .then(_ => setAuthenticated(true))
      .catch(_ => setAuthenticated(false))
  }, [])

  if (isAuthenticated == null) return <>Loading</>
  else if (isAuthenticated == false) return <LoginPage onLoginSuccess={() => setAuthenticated(true)}/>
  return <App />
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Root/>
  </StrictMode>,
)