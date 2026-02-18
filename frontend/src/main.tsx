
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import LandingContent from './LandingContent.tsx';
import './app.css';
import { BrowserRouter, Routes, Route } from "react-router";

function Landing() {
  return ( 
	<div className="items-center justify-center text-center">
	  <LandingContent />
	</div>
  )
}

function Room() {
  return ( 
    <App />
  )
}

createRoot(document.getElementById('root')!).render(
  <BrowserRouter>
    <Routes>
      <Route index element={<Landing />} />
      <Route path = "room" element={<Room />} />
    </Routes>
  </BrowserRouter>
)
