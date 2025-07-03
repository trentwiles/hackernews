// import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'

createRoot(document.getElementById('root')!).render(

  // STRICT MODE IS DISABLED TO STOP DOUBLE REQUESTS
  // https://chatgpt.com/share/68662232-cbf4-800a-8a3e-92aca82124c7

  // <StrictMode>
    <App />
  // </StrictMode>,
)
