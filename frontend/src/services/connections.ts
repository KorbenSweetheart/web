/* ============================================
   CONNECTIONS SERVICE
   ============================================
   Talks to the real /connections endpoints.
   Both list endpoints return just ids, so we
   build full profiles with fetchFullProfile
   (same pattern as recommendations).
   ============================================ */

import { apiGet, apiPost } from './api';
import { fetchFullProfile } from './users';
import type { UserProfile } from './mockUsers';

// GET /connections → accepted connections (the "Connected" tab).
// Returns a list of { id }, so we build each full profile.
export async function getConnections(): Promise<UserProfile[]> {
  const data = await apiGet('/connections');
  const ids: number[] = (data ?? []).map((c: { id: number }) => c.id);
  return Promise.all(ids.map((id) => fetchFullProfile(id)));
}

// GET /connections/requests → pending requests received (the "Received" tab).
export async function getConnectionRequests(): Promise<UserProfile[]> {
  const data = await apiGet('/connections/requests');
  const ids: number[] = (data ?? []).map((c: { id: number }) => c.id);
  return Promise.all(ids.map((id) => fetchFullProfile(id)));
}

// POST /connections → send a connection request to a user.
export async function sendConnectionRequest(toUserId: number) {
  return apiPost('/connections', { to_user_id: toUserId });
}

// PATCH /connections/:id → accept or decline a received request.
export async function respondToConnection(id: number, status: 'accepted' | 'declined') {
  const res = await fetch(`/connections/${id}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
    body: JSON.stringify({ status }),
  });
  if (!res.ok) throw new Error('Failed to update connection');
  return res.json();
}

// DELETE /connections/:id → remove a connection.
export async function deleteConnection(id: number) {
  const res = await fetch(`/connections/${id}`, {
    method: 'DELETE',
    headers: {
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
  });
  if (!res.ok) throw new Error('Failed to delete connection');
}