/* ============================================
   PROFILE SERVICE
   ============================================
   Sports list + save profile.

   The sports list is provisional (frontend mock)
   until Iván loads the real activities into the
   DB. When he does, this list gets replaced by a
   fetch to his endpoint.

   MOCK saves just log the profile for now.
   Set MOCK = false once /me/profile is live.
   ============================================ */

import type { Profile } from '../types';

const MOCK = true;

// Provisional sports list (id + name).
// Matches the shape of Iván's Activity table.
export const SPORTS = [
  { id: 1,  name: 'Running' },
  { id: 2,  name: 'Gym / Weightlifting' },
  { id: 3,  name: 'Cycling' },
  { id: 4,  name: 'Football' },
  { id: 5,  name: 'Padel' },
  { id: 6,  name: 'Tennis' },
  { id: 7,  name: 'Swimming' },
  { id: 8,  name: 'CrossFit' },
  { id: 9,  name: 'Yoga' },
  { id: 10, name: 'Basketball' },
  { id: 11, name: 'Climbing' },
  { id: 12, name: 'Boxing' },
];

// Training modes (Someone Special left out for now, per Iván)
export const MODES = [
  { id: 1, title: 'Silent',          desc: 'Train side by side, no chit-chat' },
  { id: 2, title: 'Social',          desc: 'Make it fun and social' },
  { id: 3, title: 'Open to anything', desc: 'Whatever works' },
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
  if (MOCK) {
    console.log('Saving profile (mock):', profile);
    return { success: true };
  }

  const res = await fetch('/me/profile', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${localStorage.getItem('token')}`,
    },
    body: JSON.stringify(profile),
  });

  if (!res.ok) throw new Error('Failed to save profile');
  return res.json();
}