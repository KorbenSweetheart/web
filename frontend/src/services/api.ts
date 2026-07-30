/* ============================================
   API HELPER
   ============================================
   One place for authenticated requests.
   Automatically attaches the token so we don't
   repeat it in every call.

   Usage: const data = await apiGet('/users/1');
   ============================================ */

async function apiGet(path: string) {
  const token = localStorage.getItem('token');

  const res = await fetch(path, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
  });

  if (res.status === 404) {
    throw new Error('Not found');
  }
  if (!res.ok) {
    throw new Error(`Request failed: ${res.status}`);
  }

  return res.json();
}

export { apiGet };