import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import Endpoints from './pages/Endpoints'
import EndpointDetail from './pages/EndpointDetail'
import Events from './pages/Events'
import EventDetail from './pages/EventDetail'
import DeadLetters from './pages/DeadLetters'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="endpoints" element={<Endpoints />} />
          <Route path="endpoints/:id" element={<EndpointDetail />} />
          <Route path="events" element={<Events />} />
          <Route path="events/:id" element={<EventDetail />} />
          <Route path="dead-letters" element={<DeadLetters />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
