/* ============================================
   SHARED TYPES
   ============================================
   NOTE: this file is named "mockUsers" for
   historical reasons, but it no longer holds
   any mock data — only the shared types below.
   Discover and Connections now use real backend
   data. TODO: rename this file to "types.ts"
   (low priority — it's imported in many places).
   ============================================ */

export interface Activity {
  id: number;
  title: string;
  experience: number;      // 1-5
  interest_level: number;  // 1-5
}

export interface UserProfile {
  id: number;
  name: string;
  age: number;
  picture_url: string;
  bio: string;
  max_radius: number;
  interaction_mode: number;   // 1 Open, 2 Social, 3 Silent
  activities: Activity[];
  lat: number;
  lon: number;
  is_online: boolean;
  // Not from backend — the match score comes from /recommendations later.
  // For now we mock it so the card can show the ring.
  match_score: number;        // 0-100
}