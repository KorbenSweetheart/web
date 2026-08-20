# TODO — Frontend (Pulse / match-me)

Frontend task list. Updated as Ivan progresses on the backend.

---

## Done

- [x] Real auth: register + login against the backend
- [x] `PrivateRoute` protects `/app/*` (now validates the session against the backend, since the cookie is HttpOnly and JS can't read it)
- [x] Layout: sidebar (desktop) + header/bottom-nav (mobile), with logout
- [x] Design system in `index.css` (accent colors, typography, spacing, `.sport-tag`)
- [x] Profile setup form (`ProfileSetupPage`): name, age, photo, bio, sports w/ level, mode, distance
- [x] Profile SAVE connected (`profile.ts`): PATCH `/me/profile`
- [x] Profile READ connected (`users.ts` `getMyProfile`): combines `/me` + `/me/profile` + `/me/bio`
- [x] `ProfilePage` (view my profile): real data + empty state + "Complete my profile" button
- [x] Edit flow: Edit -> pre-filled form -> save -> back to `/app/profile`
- [x] `UserCard` reusable: `variant` prop swaps footer buttons (Connect/Dismiss, Accept/Decline, Message/Remove)
- [x] `UserCard`: fallback for unknown `interaction_mode` (backend seed can emit mode 4 from paused Dating — see note below)
- [x] `ProfilePanel` (detail side panel) built and used in Discover
- [x] "Complete profile first" flow: `isProfileComplete()` + `RequireProfile` guard on Discover only. Nav menu always visible.
- [x] Redirect to login on 401 (expired/invalid token) in `api.ts`
- [x] Discover connected to REAL `/recommendations` (`getRecommendations` -> `fetchFullProfile` per id). Confirmed working with real seed users.
- [x] Location capture on profile save (browser geolocation -> lat/lon in PATCH)
- [x] Location refresh on Discover open (`updateMyLocation` before fetching recommendations)
- [x] age migration: switched from birth_date to age (number) everywhere, per Ivan
- [x] Shared `fetchFullProfile(id)` in `users.ts` (id -> full UserProfile), reused by recommendations + connections
- [x] Connections page connected to REAL endpoints (`connections.ts`):
      - `getConnectionRequests()` -> Received tab, `getConnections()` -> Connected tab
      - Accept -> PATCH accepted (moves person Received -> Connected)
      - Decline -> PATCH declined
      - Remove -> DELETE (on Connected tab)
- [x] Discover "Connect" sends a REAL request (`sendConnectionRequest` -> POST /connections)
- [x] Profile picture UPLOAD (real file upload, `POST /me/picture` / `DELETE /me/picture`)
- [x] Removed "Sent" tab (not in task requirements, no backend endpoint for it)
- [x] Cleaned up mocks: deleted `mockConnections.ts`, removed `MOCK_USERS` data (kept `mockUsers.ts` only for the `UserProfile` type)
- [x] Bigger/clearer Connections tabs (css)

---

## Done — Auth cookie migration (completed today)

- [x] Migrated auth from localStorage `Bearer` token to HttpOnly cookies
      - `api.ts`: removed token reading, added `credentials: 'include'` on every request
      - Added `apiPatch` + `apiDelete` helpers (were missing — that's why raw fetches existed)
      - `auth.ts`: login/register/logout now rely on the cookie; logout clears it via backend
      - `LoginPage.tsx`: stopped saving token to localStorage
      - `PrivateRoute.tsx`: now asks the backend "who am I?" (getMyProfile) instead of reading localStorage
- [x] Fixed all raw `fetch` calls that still sent the old `Bearer` token by hand
      - `connections.ts`: `respondToConnection`, `deleteConnection` -> use apiPatch/apiDelete
      - `users.ts`: `updateMyLocation`, `uploadProfilePicture` (kept raw for FormData + credentials), `deleteProfilePicture`
      - `profile.ts`: `saveProfile` -> apiPatch
- [x] Fixed profile save failing when photo removed: backend rejects empty `picture_url` (url validation, 400). Now we only send `picture_url` when it's non-empty; photos are managed via the `/me/picture` endpoints.
- [x] Backend CORS updated by Ivan for credentials (AllowCredentials + exact origin)
- [x] Tested: login, refresh-stays-logged-in, logout all working with the cookie

---

## Pending — my work, not blocked

- [ ] CHAT (frontend) — biggest remaining piece. Ivan's backend is merged (PR #15) with REST + WebSocket. Build order:
      1. `services/chat.ts` — REST: `POST /chats/direct`, `GET /chats`, `GET /chats/:id/messages` (paginated)
      2. Chat list view (left column: name, photo, last message preview, unread, online dot)
      3. Conversation view (messages, input, send) — mobile: 1 view at a time via `selectedChatId`
      4. WebSocket connection (the real-time layer) — needs the cookie on the handshake
      5. Live extras: typing indicator, read receipts, presence (online/offline)
      - See Ivan's chat docs (REST endpoints + WS wire protocol) for exact payloads.
- [ ] Side panel in Connections (like Discover's `ProfilePanel`)
      - Needs `ProfilePanel` made variant-aware (buttons currently fixed to Discover's Connect/Dismiss). Swap buttons per context like UserCard does.
- [ ] Public pages navigation: Login/Register have no way back to landing except browser back. Add logo->landing + link between login/register.
- [ ] `mockUsers.ts` is misnamed (only holds the `UserProfile` type). Rename to `types.ts` someday — low priority, touches many imports.

---

## Pending — waiting on backend or decisions (Ivan)

- [ ] DISCOVER: connected/dismissed users reappear when reopening Discover. ROOT CAUSE IS BACKEND — confirmed in Ivan's own code:
      - `service/match.go:49` — Ivan's comment: "don't forget to exclude users who was previously declined." So `/recommendations` doesn't filter yet.
      - `service/connection.go:18` — `AllConnectionRecords` commented out with "needed for recommendations".
      - NEED FROM IVAN:
        (a) `/recommendations` should exclude users I already have ANY connection with (pending/accepted/declined), not just declined — otherwise they reappear after Connect.
        (b) Dismiss has NO endpoint — backend never learns I skipped someone. Decide: add a skip/dismiss endpoint (+ exclude in recommendations), OR drop the Dismiss button so it's honest.
      - MY SIDE: `handleConnect`/`handleDismiss` currently only remove from local state (works in-session, lost on reload). Nothing to fix here until Ivan defines the design.
      - STATUS: message to Ivan drafted, NOT SENT YET. Send it.
- [ ] `interaction_mode` off-by-one in backend seed (`seed.go:379` uses `rand.Intn(4)` -> can emit mode 4, but only 3 modes exist since Dating is paused). Also `dto.go:25` validation allows `lte=4`. Told Ivan; my UserCard fallback covers it either way. If Dating comes back, add mode 4 to frontend MODES.
- [ ] Discover `match_score` still 0 (backend doesn't send a score yet; algorithm being designed with Ivan)
- [ ] Sports/modes from backend: `/activities` endpoint is active — can swap the hardcoded SPORTS list to fetch it
- [ ] Backend also blocking private endpoints if profile incomplete (security) — Ivan's call

---

## Chat — data shapes to align with Ivan (from his docs)

- REST:
  - `POST /chats/direct` { target_user_id } -> chat object (id, user_one, user_two, each with id/name/picture_url). 400 if no accepted connection.
  - `GET /chats` -> list of chat objects
  - `GET /chats/:id/messages?limit=15&last_message_id=105` -> paginated messages (id, chat_id, sender_id, content, created_at, is_viewed)
- WebSocket `/ws` (auth via cookie on handshake):
  - `chat:message` (send/receive), `chat:typing` (debounce 3s), `chat:read` (read receipts), `presence:check` -> `presence:batch`, `error`
- Design already done: 2 panels desktop (list + conversation), 1 view at a time mobile (`selectedChatId`: null=list, id=conversation). Reuse existing chat CSS classes. Mockup sent to Ivan.

---

## Match algorithm (design discussion with Ivan, ongoing)

- Hard filters already done in backend
- Score/ranking is the TODO. Agreed ideas: shared sports (weight), mode compatibility (weight, NOT filter — "Open to anything" is a wildcard), experience closeness (medium), distance-within-radius (matters more for running, less for gym/climbing), age by broad ranges (low).
- NOTE: interest_level is fixed at 3 for everyone, so it can't differentiate — don't use it in scoring.

---

## Decisions settled with Ivan

- [x] Auth: HttpOnly cookies + `credentials: 'include'` (chosen over localStorage/in-memory because the WS handshake can't carry a Bearer header)
- [x] Dating mode + Gender: PAUSED (may return later — not fully removed). Frontend has 3 modes; backend still has traces (seed + dto allow 4).
- [x] Availability/schedule: dropped
- [x] Interest level per sport: fixed at 3 (user doesn't pick)
- [x] Modes/levels mapped in frontend; Activities is a DB dictionary (`/activities`)
- [x] Guards: only Discover blocked if profile incomplete
- [x] Location: lat/lon in PATCH /me/profile + refresh on Discover open
- [x] age instead of birth_date (user types a number)
- [x] Connection status sent as string ("pending"/"accepted"/"declined")
- [x] No "Sent" tab (not required, no endpoint)

---

## Backend status note

- Recommendations active (but no exclusion filter yet — see Discover blocker). `/activities` active. Connections fully done + tested. 100 seed users across Finland. Profile picture upload (MinIO S3) done. Chat merged (REST + WebSocket, PR #15).
- Auth: cookie model live, CORS set for credentials.

---

## Reminder

snake_case everywhere: `access_token`, `picture_url`, `activity_id`, `interaction_mode`, `max_radius`, `experience_level`, `interest_level`, `lat`, `lon`, `to_user_id`, `from_user_id`, `target_user_id`, `chat_id`, `sender_id`, `is_viewed`, `last_message_id`.