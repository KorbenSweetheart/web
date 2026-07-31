/* ============================================
   MOCK USERS (temporary)
   ============================================
   Fake users with the EXACT shape the backend
   returns, so the UserCard can be built now.
   When /recommendations + the profile endpoints
   work, these get replaced by real fetches.

   Shape matches Iván's ProfileResponse DTO.
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

export const MOCK_USERS: UserProfile[] = [
  {
    id: 1,
    name: 'Marcus K.',
    age: 28,
    picture_url: 'https://i.pravatar.cc/300?img=12',
    bio: 'Early riser, love a hard morning run before work. Looking for someone to keep the pace honest.',
    max_radius: 15,
    interaction_mode: 3, // Silent
    activities: [
      { id: 1, title: 'Running', experience: 4, interest_level: 5 },
      { id: 3, title: 'Gym', experience: 3, interest_level: 3 },
    ],
    lat: 62.89,
    lon: 27.68,
    is_online: true,
    match_score: 92,
  },
  {
    id: 2,
    name: 'Sofia R.',
    age: 24,
    picture_url: 'https://i.pravatar.cc/300?img=45',
    bio: 'Climber and yoga person. Happy to chat between sets or just vibe in silence.',
    max_radius: 10,
    interaction_mode: 2, // Social
    activities: [
      { id: 11, title: 'Climbing', experience: 5, interest_level: 5 },
      { id: 9, title: 'Yoga', experience: 3, interest_level: 4 },
    ],
    lat: 62.9,
    lon: 27.65,
    is_online: false,
    match_score: 78,
  },
  {
    id: 3,
    name: 'Diego M.',
    age: 31,
    picture_url: 'https://i.pravatar.cc/300?img=33',
    bio: 'Football on weekends, gym during the week. Always up for a kickabout.',
    max_radius: 20,
    interaction_mode: 1, // Open to anything
    activities: [
      { id: 5, title: 'Football', experience: 4, interest_level: 4 },
      { id: 3, title: 'Gym', experience: 4, interest_level: 3 },
    ],
    lat: 62.88,
    lon: 27.7,
    is_online: true,
    match_score: 65,
  },
];