import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './app.css';
import LandingContent from './LandingContent';

createRoot(document.getElementById('root')!).render(
  <StrictMode>

    <div className="items-center justify-center text-center">
      <LandingContent />
    </div>

  </StrictMode>
)
