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
  const [profile, bio] = await Promise.all([
    apiGet('/me/profile'),
    apiGet('/me/bio'),
  ]);

  // Merge both responses into one object matching UserProfile shape
  return {
    id: profile.id,
    name: profile.name ?? '',
    age: profile.age ?? 0,
    picture_url: profile.picture_url ?? '',
    bio: profile.bio ?? '',
    max_radius: bio.max_radius ?? 0,
    interaction_mode: bio.interaction_mode ?? 1,
    activities: bio.activities ?? [],
    lat: profile.lat ?? 0,
    lon: profile.lon ?? 0,
    is_online: profile.is_online ?? false,
    match_score: 0, // not applicable to your own profile
  };
}