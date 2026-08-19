import { useState, useEffect, type ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { getMyProfile, isProfileComplete } from '../services/users';

// Guard that sits between the private routes and their content.
// It checks whether the user's profile is complete.
// - Still loading  → show nothing yet (avoids a wrong redirect flash)
// - Incomplete     → send to /app/profile-setup (unless already there)
// - Complete       → let the page render
export default function RequireProfile({ children }: { children: ReactNode }) {
  const [status, setStatus] = useState<'loading' | 'complete' | 'incomplete'>('loading');
  const location = useLocation();

  useEffect(() => {
    getMyProfile()
      .then((profile) => {
        setStatus(isProfileComplete(profile) ? 'complete' : 'incomplete');
      })
      .catch(() => {
        // If we can't load the profile, treat it as incomplete
        // (safer: pushes them to setup rather than into a broken page).
        setStatus('incomplete');
      });
  }, []);

  // While we're asking the backend, don't render anything yet.
  if (status === 'loading') {
    return null;
  }

  // Incomplete profile: force them to setup.
  // Exception: if they're ALREADY on the setup page, let them stay
  // (otherwise we'd redirect setup → setup forever = infinite loop).
  const onSetupPage = location.pathname === '/app/profile-setup';
  if (status === 'incomplete' && !onSetupPage) {
    return <Navigate to="/app/profile-setup" replace />;
  }

  // All good: show the page.
  return <>{children}</>;
}