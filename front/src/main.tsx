import { StrictMode, useEffect, useState } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { api } from './services/api.ts'
import { LoginPage } from './components/LoginPage.tsx'



const Root = () => {
  const [isAuthenticated, setAuthenticated] = useState<boolean|null>(null)

  useEffect(() => {
    const checkAuth = async () => {
      try {
        await api.me();
        setAuthenticated(true);
      } catch (e) {
        console.error("Auth check failed:", e); 
        setAuthenticated(false);
      }
    };

    checkAuth();
  }, [])

  if (isAuthenticated === null) {
    return (
        <div className="flex h-screen w-full items-center justify-center">
            <div className="text-xl font-semibold text-gray-500">Loading application...</div>
        </div>
    )
  }

  if (isAuthenticated === false) {
     return <LoginPage onLoginSuccess={() => setAuthenticated(true)}/>
  }
  
  return <App />
}
createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Root/>
  </StrictMode>,
)