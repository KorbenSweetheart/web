/* ============================================
   SHARED LABELS
   ============================================
   Single source of truth for the human-readable
   labels used across profile, panel and cards.
   Change a label here and every view updates.
   ============================================ */

// Experience level (1-5) → readable label.
export const EXP_LABELS: Record<number, string> = {
  1: 'Beginner',
  2: 'Active Novice',
  3: 'Intermediate',
  4: 'Advanced',
  5: 'Professional',
};

// Interest level (1-5) → how keen the user is on the sport.
export const INTEREST_LABELS: Record<number, string> = {
  1: 'Curious',
  2: 'Casual',
  3: 'Into it',
  4: 'Keen',
  5: 'Obsessed',
};