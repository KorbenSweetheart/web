/* ============================================
   API HELPER
   ============================================
   One place for authenticated requests.
   The browser holds the JWT in an HttpOnly cookie
   and attaches it automatically — we no longer
   read or send the token by hand.

   Usage:
     const data = await apiGet('/users/1');
     await apiPost('/auth/logout');
   ============================================ */

// No token to read anymore. We only declare that our POST bodies
// are JSON. The cookie travels on its own (see credentials below).
const jsonHeaders = {
  'Content-Type': 'application/json',
};

// Called when the backend says we're not authenticated (401).
// The cookie is HttpOnly, so there's nothing for JS to clear —
// we just send the user back to login. window.location is used
// (not navigate) because this file isn't a React component.
function handleUnauthorized() {
  window.location.href = '/login';
}

// GET: fetch data from the backend.
async function apiGet(path: string) {
  const res = await fetch(path, {
    headers: jsonHeaders,
    credentials: 'include', // ← send the auth cookie with the request
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
    headers: jsonHeaders,
    credentials: 'include', // ← same here
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

  if (res.status === 204) {
    return null;
  }
  return res.json();
}

// PATCH: update part of a resource (e.g. accept/decline a connection).
async function apiPatch(path: string, body?: unknown) {
  const res = await fetch(path, {
    method: 'PATCH',
    headers: jsonHeaders,
    credentials: 'include',
    body: body ? JSON.stringify(body) : undefined,
  });

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

  if (res.status === 204) {
    return null;
  }
  return res.json();
}

// DELETE: remove a resource (e.g. delete a connection).
async function apiDelete(path: string) {
  const res = await fetch(path, {
    method: 'DELETE',
    headers: jsonHeaders,
    credentials: 'include',
    body: undefined,
  });

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

  if (res.status === 204) {
    return null;
  }
  return res.json();
}

export { apiGet, apiPost, apiPatch, apiDelete };