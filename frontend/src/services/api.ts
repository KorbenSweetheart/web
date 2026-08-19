/* ============================================
   API HELPER
   ============================================
   One place for authenticated requests.
   Automatically attaches the token so we don't
   repeat it in every call.

   Usage:
     const data = await apiGet('/users/1');
     await apiPost('/auth/logout');
   ============================================ */

// Reads the token once, builds the auth headers.
// Shared by every request helper below.
function authHeaders() {
  const token = localStorage.getItem('token');
  return {
    'Content-Type': 'application/json',
    Authorization: `Bearer ${token}`,
  };
}

// Called when the backend says our token is no longer valid (401).
// The token is expired/invalid, so we clear it and send the user
// back to login. window.location is used (not navigate) because this
// file isn't a React component and can't use router hooks.
function handleUnauthorized() {
  localStorage.removeItem('token');
  window.location.href = '/login';
}

// GET: fetch data from the backend.
async function apiGet(path: string) {
  const res = await fetch(path, {
    headers: authHeaders(),
  });

  // Token expired or invalid → force logout + redirect to login.
  if (res.status === 401) {
    handleUnauthorized();
    throw new Error('Unauthorized');
  }
  if (res.status === 404) {
    throw new Error('Not found');
  }
  if (!res.ok) {
    throw new Error(`Request failed: ${res.status}`);
  }

  return res.json();
}

// POST: send an action to the backend (e.g. logout).
// `body` is optional — logout doesn't need one.
async function apiPost(path: string, body?: unknown) {
  const res = await fetch(path, {
    method: 'POST',
    headers: authHeaders(),
    body: body ? JSON.stringify(body) : undefined,
  });

  // Token expired or invalid → force logout + redirect to login.
  if (res.status === 401) {
    handleUnauthorized();
    throw new Error('Unauthorized');
  }
  if (res.status === 404) {
    throw new Error('Not found');
  }
  if (!res.ok) {
    throw new Error(`Request failed: ${res.status}`);
  }

  // Some endpoints (like logout) reply with no content (204).
  // Trying to parse JSON there would crash, so we guard for it.
  if (res.status === 204) {
    return null;
  }
  return res.json();
}

export { apiGet, apiPost };