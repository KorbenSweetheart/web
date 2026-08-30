/* ============================================
   CONNECTIONS SERVICE
   ============================================
   Talks to the real /connections endpoints.
   Both list endpoints return just ids, so we
   build full profiles with fetchFullProfile
   (same pattern as recommendations).
   ============================================ */

import { apiGet, apiPost, apiPatch, apiDelete } from './api';
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

// Quick count of received pending requests without fetching full profiles.
export async function getConnectionRequestsCount(): Promise<number> {
  const data = await apiGet('/connections/requests');
  return Array.isArray(data) ? data.length : 0;
}


// POST /connections → send a connection request to a user.
export async function sendConnectionRequest(toUserId: number) {
  return apiPost('/connections', { to_user_id: toUserId });
}

// PATCH /connections/:id → accept or decline a received request.
export async function respondToConnection(id: number, status: 'accepted' | 'declined') {
  return apiPatch(`/connections/${id}`, { status });
}

// DELETE /connections/:id → remove a connection.
export async function deleteConnection(id: number) {
  return apiDelete(`/connections/${id}`);
}