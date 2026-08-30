/* ============================================
   API HELPER
   ============================================
   One place for authenticated requests.
   The browser holds the JWT in an HttpOnly cookie
   and attaches it automatically — we no longer
   read or send the token by hand.

   On a 401, we transparently try to refresh the
   session once (POST /auth/refresh) and retry the
   original request. Concurrent 401s share a single
   refresh call (see refreshSession below).

   Usage:
     const data = await apiGet('/users/1');
     await apiPost('/auth/logout');
   ============================================ */

// No token to read anymore. We only declare that our POST bodies
// are JSON. The cookie travels on its own (see credentials below).
const jsonHeaders = {
  'Content-Type': 'application/json',
};

// Called when refresh fails — the session is genuinely gone.
// The cookie is HttpOnly, so there's nothing for JS to clear —
// we just send the user back to login. window.location is used
// (not navigate) because this file isn't a React component.
function handleUnauthorized() {
  window.location.href = '/login';
}

// Holds the in-flight refresh call, if any. The first request to hit
// a 401 starts the refresh and stores its promise here; any other
// request that hits a 401 while it's running awaits THIS same promise
// instead of firing its own. That's the "single refresh, queued
// requests" behaviour Ivan asked for. Reset back to null when done.
let refreshPromise: Promise<boolean> | null = null;

// Fire POST /auth/refresh exactly once, even under concurrent 401s.
// Returns true if the session was renewed, false if it's truly expired.
// The refresh_token cookie is sent automatically (credentials: include);
// on success the backend sets fresh auth cookies on the response.
function refreshSession(): Promise<boolean> {
  if (!refreshPromise) {
    refreshPromise = fetch('/auth/refresh', {
      method: 'POST',
      headers: jsonHeaders,
      credentials: 'include',
    })
      .then((res) => res.ok)
      .catch(() => false)
      .finally(() => {
        // Clear the gate so the NEXT expiry can trigger a new refresh.
        refreshPromise = null;
      });
  }
  return refreshPromise;
}

// The single fetch path every helper goes through. Handles the 401 →
// refresh → retry flow in one place, so we don't repeat it four times.
//
// `retry` guards against loops: the retried request runs with retry=false,
// so if it 401s again we give up and go to login instead of refreshing
// forever.
async function request(
  path: string,
  options: RequestInit = {},
  retry = true,
): Promise<unknown> {
  const res = await fetch(path, {
    ...options,
    headers: jsonHeaders,
    credentials: 'include', // ← send the auth cookie with every request
  });

  // Access token expired or invalid.
  if (res.status === 401) {
    // Never try to refresh the refresh call itself — that would loop.
    // If refresh returns 401, the session is genuinely over.
    if (retry && path !== '/auth/refresh') {
      const renewed = await refreshSession();
      if (renewed) {
        // Got fresh cookies — replay the original request once.
        return request(path, options, false);
      }
    }
    // Refresh failed (or we already retried) → the session is gone.
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

// GET: fetch data from the backend.
async function apiGet(path: string) {
  return request(path, { method: 'GET' });
}

// POST: send an action to the backend. `body` is optional.
async function apiPost(path: string, body?: unknown) {
  return request(path, {
    method: 'POST',
    body: body ? JSON.stringify(body) : undefined,
  });
}

// PATCH: update part of a resource.
async function apiPatch(path: string, body?: unknown) {
  return request(path, {
    method: 'PATCH',
    body: body ? JSON.stringify(body) : undefined,
  });
}

// DELETE: remove a resource.
async function apiDelete(path: string) {
  return request(path, { method: 'DELETE' });
}

export { apiGet, apiPost, apiPatch, apiDelete };