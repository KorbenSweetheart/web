import { useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { getMyProfile } from '../services/users';

// Three states while we ask the backend "is this cookie valid?":
//   'checking' → request in flight, we don't know yet
//   'ok'       → backend answered fine, the cookie is valid
//   'no'       → backend said 401 (or errored), no valid session
type AuthState = 'checking' | 'ok' | 'no';

export default function PrivateRoute({ children }: { children: ReactNode }) {
  const [status, setStatus] = useState<AuthState>('checking');

  useEffect(() => {
    // We can't read the HttpOnly cookie from JS, so the only way to
    // know if we're logged in is to make an authenticated request and
    // see if it succeeds. getMyProfile() hits a protected endpoint,
    // so it doubles as our "am I logged in?" check.
    getMyProfile()
      .then(() => setStatus('ok'))
      .catch(() => setStatus('no'));
  }, []);

  // Still waiting on the backend — show a placeholder so the page
  // doesn't flash or bounce before we actually know the answer.
  if (status === 'checking') {
    return <div className="route-loading">Loading…</div>;
  }

  if (status === 'no') {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
}