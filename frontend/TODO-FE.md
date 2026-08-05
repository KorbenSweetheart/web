# TODO — Frontend (Pulse / match-me)

Frontend task list. Updated as Iván progresses on the backend.

---

## ✅ Done

- Real auth: register + login work against the backend (real JWT in localStorage).
- `PrivateRoute` protects `/app/*`.
- Full layout: sidebar (desktop) + header/bottom-nav (mobile), with logout button.
- Design system in `index.css`: 3 accent colors, typography classes, spacing utils, `.sport-tag`.
- **Profile setup form** (`ProfileSetupPage`): name, date of birth, photo URL, bio, sports w/ level, mode, distance.
- **Real profile SAVE connected** (`profile.ts`, MOCK=false): PATCH `/me/profile`, translates our field names to backend shape (`interaction_mode`, activity `id`/`experience_level`).
- **Profile READ connected** (`users.ts` `getMyProfile`): combines `/me` + `/me/profile` + `/me/bio` into one object.
- **ProfilePage** (view my profile): shows real data, empty state when blank.
- **Edit flow works**: Edit button → pre-filled form → save → back to `/app/profile`.
- Form pre-fills from existing profile (useEffect + getMyProfile). First-timers get empty form.
- Redirect after save: editing → `/app/profile`; first-time → `/app/discover`.
- **Discover page** (mock data): `DiscoverPage`, `UserCard`, `ProfilePanel` (master-detail; side panel desktop, fullscreen mobile).
- Routes exist for `/app/profile` (view) and `/app/profile-setup` (fill/edit).
- Backend profile bug (InteractionMode) — fixed by Iván. All read endpoints work.
- **birth_date migration DONE**: `age` → `birth_date` (`YYYY-MM-DD`) across `types/index.ts`, `services/profile.ts`, `services/users.ts`, `services/mockUsers.ts`, `ProfileSetupPage.tsx`, `ProfilePage.tsx`, `UserCard.tsx`. Age shown to the user is computed on the fly with `calculateAge()` (display only, never sent to backend).
- **Connections page (mock)**: `ConnectionsPage.tsx` + `ConnectionsPage.css`, three tabs (Received / Sent / Connected), backed by `services/mockConnections.ts` (splits `MOCK_USERS`). Route `/app/connections` wired in `routes.tsx`.
- **`UserCard` reusable across pages**: added `variant` prop (`discover` | `received` | `sent` | `connected`) that swaps footer buttons (Connect/Dismiss, Accept/Decline, Cancel, Message) without duplicating the card.

---

## ⚠️ MOCKS still in place (swap for real data later)

### Discover uses mock users
- `services/mockUsers.ts` → `MOCK_USERS` (Marcus, Sofia, Diego), built with EXACT backend DTO shape (now with `birth_date`).
- `DiscoverPage.tsx` imports MOCK_USERS instead of fetching.
- Needs `GET /recommendations` (list of IDs) — NOT active yet.
- Flow when ready: `/recommendations` → per id `/users/:id` + `/users/:id/bio` → combine → cards.
- `match_score` is faked in the mock; comes from `/recommendations` later.
- Connect/Dismiss buttons only do `console.log` + remove from local state. Later: call backend.

### Connections uses mock data
- `services/mockConnections.ts` splits `MOCK_USERS` into `RECEIVED` / `SENT` / `CONNECTED` — just for layout, not realistic.
- Needs `GET /connections` (NOT active yet on backend — commented in router.go). Real shape likely includes a `status` field instead of three separate arrays; adjust `ConnectionsPage.tsx` once known.
- Accept/Decline/Cancel/Message buttons only do `console.log` + local state changes. Later: call backend.
- Message button should eventually `navigate('/app/chats/:id')` once Chat exists.

---

## ⏳ Pending (my code, waiting on backend or decisions)

### Auth: migrate to HttpOnly cookies + refresh flow ← NEW, Iván confirmed this direction
- Iván is moving away from body/localStorage token toward **cookie-only** auth:
  - Login sets `access_token` cookie (`Path: /`) and `refresh_token` cookie (`Path: /auth/refresh`), both HttpOnly.
  - On 401 from any private route → frontend calls `POST /auth/refresh` (browser auto-attaches the refresh cookie) → backend issues new cookies → retry original request.
- This means rewriting `services/api.ts`:
  - Remove manual `Authorization: Bearer` + `localStorage.getItem('token')`.
  - Add `credentials: 'include'` to every fetch.
  - Add 401 → refresh → retry logic (probably wrapped once in `apiGet`/`apiPost` so it's automatic everywhere).
- Also touches `services/auth.ts` (no more storing token on login) and `logout()` (no more `localStorage.removeItem('token')`).
- **Not started yet.** Bigger than originally scoped as "token expiry" below — supersedes it.
- Waiting on: confirmation that `/auth/refresh` is actually active on the backend before wiring the retry logic.

### "Complete profile first" flow ← NOT built yet
- Task: user must complete profile before Discover/Chats/Connections.
- Both routes exist but there's NO logic deciding which to show.
- Need a signal: `profile_completed` field on `/me`, OR infer (empty `activities` = incomplete).
- Then in `PrivateRoute`: incomplete → redirect to `/app/profile-setup`; hide menu while incomplete.
- Backend should also block those endpoints if profile incomplete (real security).

### Sports/modes from backend
- `profile.ts` SPORTS/MODES hardcoded (correct seed IDs). Swap to fetch when `/activities` exists.
- Modes/levels agreed to stay mapped in frontend.

### /me/profile doesn't return name/picture
- `/me/profile` returns only id/bio (+ now birth_date). name + picture_url come from `/me`.
- That's why getMyProfile combines all three. Noted in case it confuses later.

### Profile picture upload
- Iván flagged: current picture_url (pasted URL) is "100% not the right behaviour" long-term — needs real upload/storage.
- No frontend work yet; wait for backend direction on how uploads will be handled.

---

## 🆕 New work (build)

- **Connections real data**: swap `mockConnections.ts` for `GET /connections` once active; wire Accept/Decline/Cancel to real endpoints.
- **Chat**: needs connections + WebSocket (real-time, no polling). Big piece: chat list, messages, unread badge, typing indicator, online status. Not started.

---

## 💬 Decisions settled with Iván

- Dating mode + Gender: DROPPED (training-partners only). Revisit Dating at end if time.
- Availability/schedule: DROPPED (over-engineers recommendations).
- Interest level per sport: FIXED at 3 in form (backend supports it, user doesn't pick).
- Only Activities becomes a DB dictionary (`/activities`); modes/levels mapped in frontend.
- birth_date instead of age: ✅ DONE on frontend. Format `YYYY-MM-DD` (ISO 8601), confirmed by Iván.
- Enums (InteractionMode, ExperienceLevel, InterestLevel) confirmed as backend enums, not tables — frontend mapping unchanged.

## ⚙️ Backend status note
- Iván: PostGIS added (crashed Docker doing it), GORM has been a struggle. Next focus: activities + connection-to-profile, then recommendations. Backend may be unstable for a bit — if requests fail weirdly, might be him.
- Auth moving to cookie + refresh model (see pending item above) — this is a confirmed direction, not just a note.

---

## 🔤 Reminder
Everything in snake_case: `access_token`, `picture_url`, `activity_id`, `interaction_mode`, `max_radius`, `experience_level`, `interest_level`, `profile_completed`, `birth_date`.