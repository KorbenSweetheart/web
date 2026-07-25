/* ============================================
   AUTH SERVICE
   ============================================
   Single place where the frontend talks to the
   backend about authentication.

   WHY THIS FILE EXISTS:
   Pages should not call fetch() directly. They
   just call registerUser() / loginUser() and
   don't care how the request is made.
   If the API changes, only this file changes.

   MOCK MODE:
   While the backend auth endpoints are still in
   progress, MOCK returns fake responses so the
   frontend flow can be built and tested.
   Set MOCK = false once /auth/login is live.
   ============================================ */

const MOCK = true;

export async function registerUser(email: string, password: string) {
  if (MOCK) {
    return { token: 'fake-token-123', profileCompleted: false };
  }

  const res = await fetch('/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });

  if (!res.ok) throw new Error('Registration failed');
  return res.json();
}

export async function loginUser(email: string, password: string) {
  if (MOCK) {
    return { token: 'fake-token-123', profileCompleted: true };
  }

  const res = await fetch('/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });

  if (!res.ok) throw new Error('Login failed');
  return res.json();
}

export function logout() {
  localStorage.removeItem('token');
}