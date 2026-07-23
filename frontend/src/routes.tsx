import { BrowserRouter, Routes, Route } from 'react-router-dom';
import LandingPage from './pages/LandingPage';
import LoginPage from './pages/auth/LoginPage';
import RegisterPage from './pages/auth/RegisterPage';
import AppLayout from './layouts/AppLayout';

export default function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        {/*Public Routes */}
        <Route path="/" element={<LandingPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />

         {/*Private Routes - Layout */}
        <Route path="/app" element={<AppLayout />}>
          <Route path="discover" element={<div className="text-hero" style={{ color: 'var(--text-1)' }}>🔥 Discover</div>} />
          <Route path="connections" element={<div className="text-hero" style={{ color: 'var(--text-1)' }}>👥 Connections</div>} />
          <Route path="chats" element={<div className="text-hero" style={{ color: 'var(--text-1)' }}>💬 Chats</div>} />
          <Route path="profile" element={<div className="text-hero" style={{ color: 'var(--text-1)' }}>👤 Profile</div>} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}