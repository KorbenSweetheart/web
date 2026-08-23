import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import LandingPage from './pages/LandingPage';
import LoginPage from './pages/auth/LoginPage';
import RegisterPage from './pages/auth/RegisterPage';
import AppLayout from './layouts/AppLayout';
import PrivateRoute from './components/PrivateRoute';
import RequireProfile from './components/RequireProfile';
import ProfileSetupPage from './pages/ProfileSetupPage';
import DiscoverPage from './pages/DiscoverPage';
import ProfilePage from './pages/ProfilePage';
import ConnectionsPage from './pages/ConnectionsPage';
import ChatsPage from './pages/ChatsPage';


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
          <Route path="discover" element={<RequireProfile><DiscoverPage /></RequireProfile>} />
          <Route path="connections" element={<ConnectionsPage />} />
          <Route path="chats" element={<ChatsPage />} />
          <Route path="profile" element={<ProfilePage />} />
          <Route path="profile-setup" element={<ProfileSetupPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}