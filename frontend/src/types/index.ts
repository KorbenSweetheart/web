/* ============================================
   SHARED TYPES
   ============================================
   Shape of the data, matching Iván's backend
   model (domain.go). Keep snake_case to match
   the API responses exactly.
   ============================================ */

// One sport the user practices, with their level in it
export interface ProfileActivity {
  activity_id: number;      // which sport (from the activities list)
  experience: number;       // 1-5: Beginner → Professional
  interest_level: number;   // 1-5: how keen they are on it
}

// The full profile the user fills in
export interface Profile {
  name: string;
  birth_date: string;
  bio: string;
  picture_url: string;
  interaction_mode_id: number;   // Silent / Social / Open to anything
  activities: ProfileActivity[]; // their sports + levels
  max_radius: number;            // km, for location matching
}