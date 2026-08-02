# TODO — Frontend (Pulse / match-me)

Frontend task list. Updated as Iván progresses on the backend.
Written in English so Iván can read it too.

---

## ✅ Done

- Real auth connected: register + login work against the backend (real JWT).
- Register → sends to `/login` (no token returned, on purpose).
- Login → stores the `access_token` and enters the app.
- `PrivateRoute` protects `/app/*` (redirects to `/login` without a token).
- Vite proxy to `localhost:8080`.
- Full layout: sidebar (desktop) + header and bottom-nav (mobile), with logout.
- Profile setup form (ProfileSetupPage): name, date of birth, photo URL, bio, sports with level, mode, distance.
- Sport and mode IDs aligned with the backend seed.
- `services/api.ts` — helper that auto-attaches the token to requests.
- `services/users.ts` — `getMyProfile()` combines `/me/profile` + `/me/bio` into one object.
- **ProfilePage** (view my profile) — connected to REAL data via `/me`. Loads name, photo, sports, mode, distance. Shows an empty state when the profile has no data yet.
- **Discover page** built: `DiscoverPage`, `UserCard`, `ProfilePanel` (master-detail: card list + side panel on desktop, fullscreen on mobile).
- Routes for both `/app/profile` (view/edit) and `/app/profile-setup` (fill in) now coexist.
- Backend bug (`InteractionMode: unsupported relations`) — **fixed by Iván**, confirmed `/me`, `/me/profile`, `/me/bio`, `/users/:id` all work now.

---

## ⚠️ MOCKS still in place (must be swapped for real data later)

These parts work visually but use FAKE data. They need to be connected to real
endpoints once those are ready. Keeping them listed so nothing is forgotten:

### 1. Discover uses mock users
- `services/mockUsers.ts` → `MOCK_USERS` (3 fake users: Marcus, Sofia, Diego).
- `DiscoverPage.tsx` imports `MOCK_USERS` instead of fetching real data.
- **The mock was built with the EXACT backend shape** (matches `ProfileResponse` DTO),
  so swapping to real data should be smooth.
- **To connect for real, still need:**
  - `GET /recommendations` (returns a list of IDs) — NOT active yet in the backend.
  - Then: fetch `/recommendations` -> for each id fetch `/users/:id` + `/users/:id/bio` -> combine -> render cards.
  - The `match_score` in the mock is fake — it will come from `/recommendations` later.
- Connect / Dismiss buttons currently just `console.log` + remove the card from local state.
  Later they must call the backend (send connection request / dismiss).

### 2. Saving the profile is still mocked
- `services/profile.ts` → `MOCK = true`. `saveProfile` only does a `console.log`.
- **Before turning MOCK off:** the backend `UpdateProfile` endpoint
  (`POST /users/:id/profile`) is still COMMENTED OUT in main. Iván is working
  around the profile area (branch `ivan/profile-update`).
- Once it's active: read the backend to confirm method (PUT/POST) and exact shape,
  then set `MOCK = false`. The fetch is already written, just needs turning on.
- Must accept `activities: [{ activity_id, experience, interest_level }]`.

---

## ⏳ Pending to connect (changes in MY code, waiting on backend)

### "Complete profile first" flow
- Task requires: users must complete their profile before seeing recommendations
  or connecting. No access to Discover/Chats/Connections until the profile is done.
- Frontend part (mine, once the signal exists):
  - In `PrivateRoute`: if profile not complete -> redirect to `/app/profile-setup`.
  - Hide the menu (sidebar / bottom-nav) while the profile is incomplete.
  - Once complete -> normal access; `/app/profile` becomes view/edit.
- **Needs from Iván (open question):** a way to know if a profile is complete.
  Either a `profile_completed` field on `/me`, or infer it (empty `activities` = not done).
- Backend should ALSO block those endpoints if the profile isn't complete (real security,
  not just the frontend hiding things).

### Edit profile
- `ProfilePage` has an **Edit button that is disabled** ("Coming soon").
- Enable it once profile saving/editing works. It should reuse the ProfileSetup form,
  pre-filled with current data.

### Redirect after login by profile status
- `LoginPage.tsx` currently sends everyone to `/app/profile`.
- Once the "profile complete" signal exists:
  - complete -> `/app/discover`
  - incomplete -> `/app/profile-setup`

### Token expiry handling (401)
- When any request returns 401 (expired token), the app should auto clear the token
  and redirect to `/login`, instead of showing a confusing error.
- Add this to `services/api.ts`. Right now an expired token just shows
  "Could not load your profile".

### Sports / modes lists from the backend
- `services/profile.ts` → `SPORTS` and `MODES` are hardcoded (with correct seed IDs).
- When Iván builds `/activities` (and if he exposes modes), replace the hardcoded
  lists with a fetch. Modes/levels agreed to stay mapped in frontend.

### Real logout
- `/auth/logout` exists on the backend (still commented in main router though).
- When active, `logout()` can call it in addition to clearing the local token.

---

## 🆕 New work (build, not just connect)

### Connections
- Needs `GET /connections` (list of IDs) — Iván working on it.
- Reuse the card, different buttons per state (pending / connected -> chat, disconnect).

### Chat
- Needs connections + WebSocket (real-time, no polling per the task).
- Big separate piece: chat list, message view, unread indicator, typing indicator,
  online/offline status.

---

## 💬 Decisions settled with Iván

- **Dating mode + Gender:** DROPPED. Pulse is training-partners only (simpler). Can revisit Dating at the end if there's time.
- **Availability / schedule:** DROPPED (Iván: over-engineers the recommendation engine).
- **Interest level per sport:** kept FIXED at 3 in the form (backend supports it, but user doesn't pick it — keeps the form simple).
- **Only Activities becomes a DB dictionary** (`/activities`); modes and levels stay mapped in frontend.

## 💬 Still open with Iván

- How the frontend knows a profile is complete (field vs infer).
- Whether I take on a piece of the backend to share the load (offered, waiting on his call).

---

## 🔤 Permanent reminder

Everything in **snake_case** to match the backend:
`access_token`, `picture_url`, `activity_id`, `interaction_mode_id`, `max_radius`, `profile_completed`, `interest_level`.
