/* ============================================
   PROFILE SERVICE
   ============================================
   Sports + modes + levels lists, and save profile.

   The lists below match Iván's DB exactly (ids
   confirmed from the seeded backend). Later these
   will come from /activities and /interactionMode
   endpoints instead of being hardcoded here.

   saveProfile is still MOCK until /me/profile is
   wired up. Set MOCK = false when ready.
   ============================================ */

import type { Profile } from '../types';

const MOCK = false;

// Sports list — ids match Iván's activities table exactly.
// (Confirmed from the seeded DB, so ids line up with the backend.)
export const SPORTS = [
  { id: 1,  name: 'Running' },
  { id: 2,  name: 'Padel' },
  { id: 3,  name: 'Gym' },
  { id: 4,  name: 'Cycling' },
  { id: 5,  name: 'Football' },
  { id: 6,  name: 'Tennis' },
  { id: 7,  name: 'Swimming' },
  { id: 8,  name: 'CrossFit' },
  { id: 9,  name: 'Yoga' },
  { id: 10, name: 'Basketball' },
  { id: 11, name: 'Climbing' },
  { id: 12, name: 'Boxing' },
  { id: 13, name: 'MMA' },
  { id: 14, name: 'Aikido' },
  { id: 15, name: 'Jiu-Jitsu' },
];

// Training modes — ids match Iván's InteractionMode enum.
// "Dating" (id 4) intentionally left out — Pulse is about
export const MODES = [
  { id: 1, title: 'Open to anything', desc: 'Whatever works' },
  { id: 2, title: 'Social',           desc: 'Make it fun and social' },
  { id: 3, title: 'Silent',           desc: 'Train side by side, no chit-chat' },
];

// Experience levels (1-5, matching Iván's model)
export const LEVELS = [
  { value: 1, label: 'Beginner' },
  { value: 2, label: 'Active Novice' },
  { value: 3, label: 'Intermediate' },
  { value: 4, label: 'Advanced' },
  { value: 5, label: 'Professional' },
];

export async function saveProfile(profile: Profile) {
  // Transform our data into the exact shape Iván's backend expects.
  // Our internal names differ, so we translate here before sending:
  //   interaction_mode_id → interaction_mode
  //   activity_id         → id
  //   experience          → experience_level
  const body = {
    name: profile.name,
    picture_url: profile.picture_url,
    birth_date: profile.birth_date,
    bio: profile.bio,
    max_radius: profile.max_radius,
    interaction_mode: profile.interaction_mode_id,
    activities: profile.activities.map((a) => ({
      id: a.activity_id,
      experience_level: a.experience,
      interest_level: a.interest_level,
    })),
  };

  if (MOCK) {
    console.log('Saving profile (mock):', body);
    return { success: true };
  }

  const res = await fetch('/me/profile', {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
    body: JSON.stringify(body),
  });

  if (!res.ok) throw new Error('Failed to save profile');
  return res.json();
}