import { BrowserRouter, Route, Routes } from 'react-router-dom'
import LengthenPage from './pages/ShortenPage'
import RedirectPage from './pages/RedirectPage'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LengthenPage />} />
        <Route path="/:slug" element={<RedirectPage />} />
      </Routes>
    </BrowserRouter>
  )
}
