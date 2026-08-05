import type { UserProfile } from './mockUsers';
import { MOCK_USERS } from './mockUsers';

// Mock split into the three states, reusing MOCK_USERS as source people.
// Swap for real GET /connections (with a `status` field) when Iván activates it.
export const RECEIVED: UserProfile[] = MOCK_USERS.slice(0, 1);
export const SENT: UserProfile[] = MOCK_USERS.slice(1, 2);
export const CONNECTED: UserProfile[] = MOCK_USERS.slice(2, 3);