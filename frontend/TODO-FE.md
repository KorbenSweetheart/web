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

---

## ⚠️ MOCKS still in place (swap for real data later)

### Discover uses mock users
- `services/mockUsers.ts` → `MOCK_USERS` (Marcus, Sofia, Diego), built with EXACT backend DTO shape.
- `DiscoverPage.tsx` imports MOCK_USERS instead of fetching.
- Needs `GET /recommendations` (list of IDs) — NOT active yet.
- Flow when ready: `/recommendations` → per id `/users/:id` + `/users/:id/bio` → combine → cards.
- `match_score` is faked in the mock; comes from `/recommendations` later.
- Connect/Dismiss buttons only do `console.log` + remove from local state. Later: call backend.

---

## ⏳ Pending (my code, waiting on backend or decisions)

### "Complete profile first" flow  ← NOT built yet
- Task: user must complete profile before Discover/Chats/Connections.
- Both routes exist but there's NO logic deciding which to show.
- Need a signal: `profile_completed` field on `/me`, OR infer (empty `activities` = incomplete).
- Then in `PrivateRoute`: incomplete → redirect to `/app/profile-setup`; hide menu while incomplete.
- Backend should also block those endpoints if profile incomplete (real security).

### Date of birth vs age  ← Iván AGREED to change
- Backend currently stores `age` (number). Editing can't show a date, and age gets stale.
- Iván said YES to switching to `birth_date`. When he does: form sends a date like "1990-05-14".
- Until then: form still asks age via date picker, and on EDIT the date field comes back EMPTY
  (can't rebuild date from age), so user must re-enter it. Small annoyance, left as-is for now.

### Token expiry (401)
- When a request returns 401, app should auto-clear token + redirect to `/login`.
- Add to `services/api.ts`. Right now expired token just shows a generic error.

### Sports/modes from backend
- `profile.ts` SPORTS/MODES hardcoded (correct seed IDs). Swap to fetch when `/activities` exists.
- Modes/levels agreed to stay mapped in frontend.

### Real logout
- `/auth/logout` now exists on backend. `logout()` can call it in addition to clearing local token.

### /me/profile doesn't return name/picture
- `/me/profile` returns only id/age/bio. name + picture_url come from `/me`.
- That's why getMyProfile combines all three. Noted in case it confuses later.

---

## 🆕 New work (build)

- **Connections**: needs `GET /connections`. Reuse card, different buttons per state.
- **Chat**: needs connections + WebSocket (real-time, no polling). Big piece: chat list, messages, unread badge, typing indicator, online status.

---

## 💬 Decisions settled with Iván

- Dating mode + Gender: DROPPED (training-partners only). Revisit Dating at end if time.
- Availability/schedule: DROPPED (over-engineers recommendations).
- Interest level per sport: FIXED at 3 in form (backend supports it, user doesn't pick).
- Only Activities becomes a DB dictionary (`/activities`); modes/levels mapped in frontend.
- birth_date instead of age: Iván AGREED, will change backend when he can.

## ⚙️ Backend status note
- Iván is mid-work: switching age→birth_date, wrestling with GORM, added PostGIS for coordinates
  (crashed Docker doing it). Backend may be unstable for a bit — if requests fail weirdly, might be him.

---

## 🔤 Reminder
Everything in snake_case: `access_token`, `picture_url`, `activity_id`, `interaction_mode`, `max_radius`, `experience_level`, `interest_level`, `profile_completed`.