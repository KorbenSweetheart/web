/* ============================================
   USERS SERVICE
   ============================================
   Functions to fetch user data from the backend.
   Uses apiGet (which attaches the token).

   To build a full profile we combine several
   endpoints, as Iván's API requires.
   ============================================ */

import { request, apiGet, apiPatch, apiDelete, apiPost } from './api';
import type { UserProfile } from './mockUsers';

// Fetch the logged-in user's full profile.
// Combines /me/profile (about) + /me/bio (sports, mode, radius).
export async function getMyProfile(): Promise<UserProfile> {
  const [summary, profile, bio] = await Promise.all([
    apiGet('/me'),          // name, picture_url
    apiGet('/me/profile'),  // age, bio
    apiGet('/me/bio'),      // max_radius, mode, activities
  ]);

  return {
    id: summary.id,
    name: summary.name ?? '',
    picture_url: summary.picture_url ?? '',
    age: profile.age ?? 0,
    bio: profile.bio ?? '',
    max_radius: bio.max_radius ?? 0,
    interaction_mode: bio.interaction_mode ?? 1,
    activities: bio.activities ?? [],
    lat: profile.lat ?? 0,
    lon: profile.lon ?? 0,
    is_online: profile.is_online ?? false,
    match_score: 0,
  };
}

// Checks whether a profile has all the required fields filled in.
export function isProfileComplete(profile: UserProfile): boolean {
  // "Complete enough to match" = the fields matching actually needs.
  // Age and bio are optional (see profile setup), so they're not required
  // here — otherwise a valid saved profile would loop back to setup.
  return (
    profile.name.trim() !== '' &&        // has a name
    profile.interaction_mode > 0 &&      // picked a training mode
    profile.activities.length > 0        // picked at least one sport
  );
}

// Fetch recommendations from the backend.
// Backend returns { recommendations: [1, 2, ...] } — just a list of ids.
// For each id we combine /users/:id + /users/:id/profile + /users/:id/bio
// into one full UserProfile the card can show (same idea as getMyProfile).
// Takes a user id and builds their full UserProfile by combining
// /users/:id + /users/:id/profile + /users/:id/bio (the API requires
// combining several endpoints). Reused by recommendations and connections.
export async function fetchFullProfile(id: number): Promise<UserProfile> {
  const [summary, profile, bio] = await Promise.all([
    apiGet(`/users/${id}`),          // name, picture_url
    apiGet(`/users/${id}/profile`),  // bio (about me), age
    apiGet(`/users/${id}/bio`),      // max_radius, mode, activities
  ]);

  return {
    id: summary.id,
    name: summary.name ?? '',
    picture_url: summary.picture_url ?? '',
    age: profile.age ?? 0,
    bio: profile.bio ?? '',
    max_radius: bio.max_radius ?? 0,
    interaction_mode: bio.interaction_mode ?? 1,
    activities: bio.activities ?? [],
    lat: profile.lat ?? 0,
    lon: profile.lon ?? 0,
    is_online: profile.is_online ?? false,
    match_score: 0, // backend doesn't send a score yet
  };
}

// Fetch recommendations: backend returns
//   { recommendations: [ { user_id, score }, ... ] }
// The score comes from THIS endpoint; the rest of the profile comes from
// fetchFullProfile. So we fetch each full profile and attach its score.
export async function getRecommendations(): Promise<UserProfile[]> {
  const data = await apiGet('/recommendations');
  const recs: { user_id: number; score: number }[] = data.recommendations ?? [];
  return Promise.all(
    recs.map(async (rec) => {
      const profile = await fetchFullProfile(rec.user_id);
      return { ...profile, match_score: rec.score };
    }),
  );
}

// Sends the user's current browser location to the backend (PATCH /me/profile).
// Used before fetching recommendations so matching uses fresh coordinates.
// Resolves silently if location isn't available — we don't want to block Discover.
export async function updateMyLocation(): Promise<void> {
  return new Promise((resolve) => {
    if (!navigator.geolocation) return resolve();
    navigator.geolocation.getCurrentPosition(
      async (pos) => {
        try {
          await apiPatch('/me/profile', {
            lat: pos.coords.latitude,
            lon: pos.coords.longitude,
          });
        } catch (err) {
          console.warn('Could not update location:', err);
        }
        resolve();
      },
      () => resolve() // denied or failed → just continue
    );
  });
}

// Uploads a profile picture file (POST /me/picture).
// Sends the file as FormData. Routed through request() so it benefits from
// automatic token refresh and session recovery just like all other API calls.
export async function uploadProfilePicture(file: File): Promise<string> {
  const form = new FormData();
  form.append('picture', file); // 'picture' is the field name the backend expects

  const data = (await request('/me/picture', {
    method: 'POST',
    body: form,
  })) as { picture_url: string };

  return data.picture_url;
}

// Removes the profile picture (DELETE /me/picture), resets to default.
export async function deleteProfilePicture(): Promise<string> {
  return apiDelete('/me/picture');
}

// POST /recommendations/:id/dismiss — persists the dismissal so the
// backend stops recommending this user. Empty body; id goes in the path.
export async function dismissRecommendation(userId: number): Promise<void> {
  await apiPost(`/recommendations/${userId}/dismiss`, {});
}