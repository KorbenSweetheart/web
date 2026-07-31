import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import LandingPage from './pages/LandingPage';
import LoginPage from './pages/auth/LoginPage';
import RegisterPage from './pages/auth/RegisterPage';
import AppLayout from './layouts/AppLayout';
import PrivateRoute from './components/PrivateRoute';
import ProfileSetupPage from './pages/ProfileSetupPage';
import DiscoverPage from './pages/DiscoverPage';


export default function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        {/*Public Routes */}
        <Route path="/" element={<LandingPage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />

         {/*Private Routes - Layout */}
        <Route path="/app" element={<PrivateRoute><AppLayout /></PrivateRoute>}>
          <Route index element={<Navigate to="discover" replace />} />
          <Route path="discover" element={<DiscoverPage />} />
          <Route path="connections" element={<div className="text-hero" style={{ color: 'var(--text-1)' }}>👥 Connections</div>} />
          <Route path="chats" element={<div className="text-hero" style={{ color: 'var(--text-1)' }}>💬 Chats</div>} />
          <Route path="profile" element={<ProfileSetupPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}