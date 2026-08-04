/* ============================================
   USERS SERVICE
   ============================================
   Functions to fetch user data from the backend.
   Uses apiGet (which attaches the token).

   To build a full profile we combine several
   endpoints, as Iván's API requires.
   ============================================ */

import { apiGet } from './api';
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