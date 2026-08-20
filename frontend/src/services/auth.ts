/* ============================================
   AUTH SERVICE
   ============================================
   Single place where the frontend talks to the
   backend about authentication.

   WHY THIS FILE EXISTS:
   Pages call registerUser() / loginUser() and
   don't care how the request is made. If the API
   changes, only this file changes.

   AUTH APPROACH (HttpOnly cookie):
   Login makes the backend set an HttpOnly cookie
   with the JWT. The browser stores it and sends it
   automatically on every request that uses
   credentials: 'include'. The frontend never reads
   or stores the token itself — it's invisible to JS,
   which is what makes it safe against XSS.

   Register does NOT log the user in — they're sent
   to /login afterwards.
   ============================================ */

import { apiPost } from './api';

// Register a new account. Backend needs name, email, password.
// Returns { id, email, name, message } — no token, no cookie yet.
export async function registerUser(name: string, email: string, password: string) {
  const res = await fetch('/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include', // for consistency with every other request
    body: JSON.stringify({ name, email, password }),
  });

  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error(data.message || 'Registration failed');
  }
  return res.json();
}

// Log in. The backend replies with a Set-Cookie header holding the
// JWT. credentials: 'include' is what lets the browser actually store
// that cookie — without it, login "works" but no cookie is saved and
// every later request comes back 401.
export async function loginUser(email: string, password: string) {
  const res = await fetch('/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include', // ← lets the browser store the auth cookie
    body: JSON.stringify({ email, password }),
  });

  if (!res.ok) {
    throw new Error('Invalid email or password');
  }
  return res.json();
}

// Log out. The backend clears the cookie (it sends back a Set-Cookie
// that empties it). We can't clear an HttpOnly cookie from JS, so the
// backend call is the logout now — there's nothing left to remove
// locally.
export async function logout() {
  try {
    await apiPost('/auth/logout');
  } catch (err) {
    // Backend might be down or the session already expired.
    // We log it but don't throw, so the user isn't trapped in the app.
    console.error('Backend logout failed:', err);
  }
}