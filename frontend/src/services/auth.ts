/* ============================================
   AUTH SERVICE
   ============================================
   Single place where the frontend talks to the
   backend about authentication.

   WHY THIS FILE EXISTS:
   Pages call registerUser() / loginUser() and
   don't care how the request is made. If the API
   changes, only this file changes.

   AUTH APPROACH (header):
   Login returns access_token in the JSON body.
   We store it and send it as
   "Authorization: Bearer <token>" on protected
   requests. The backend also accepts a cookie,
   but we use the header for simplicity.

   Register does NOT return a token — the user is
   sent to /login afterwards.
   ============================================ */

// Register a new account. Backend needs name, email, password.
// Returns { id, email, name, message } — no token.
export async function registerUser(name: string, email: string, password: string) {
  const res = await fetch('/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name, email, password }),
  });

  if (!res.ok) {
    const data = await res.json().catch(() => ({}));
    throw new Error(data.message || 'Registration failed');
  }
  return res.json();
}

// Log in. Backend returns { message, access_token }.
export async function loginUser(email: string, password: string) {
  const res = await fetch('/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  });

  if (!res.ok) {
    throw new Error('Invalid email or password');
  }
  return res.json();
}

export function logout() {
  localStorage.removeItem('token');
}