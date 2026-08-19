# TODO — Frontend (Pulse / match-me)

Frontend task list. Updated as Ivan progresses on the backend.

---

## Done

- [x] Real auth: register + login against the backend (JWT in localStorage)
- [x] `PrivateRoute` protects `/app/*`
- [x] Layout: sidebar (desktop) + header/bottom-nav (mobile), with logout
- [x] Design system in `index.css` (accent colors, typography, spacing, `.sport-tag`)
- [x] Profile setup form (`ProfileSetupPage`): name, age, photo URL, bio, sports w/ level, mode, distance
- [x] Profile SAVE connected (`profile.ts`): PATCH `/me/profile`
- [x] Profile READ connected (`users.ts` `getMyProfile`): combines `/me` + `/me/profile` + `/me/bio`
- [x] `ProfilePage` (view my profile): real data + empty state + "Complete my profile" button
- [x] Edit flow: Edit -> pre-filled form -> save -> back to `/app/profile`
- [x] `UserCard` reusable: `variant` prop swaps footer buttons (Connect/Dismiss, Accept/Decline, Message/Remove)
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
- [x] Removed "Sent" tab (not in task requirements, no backend endpoint for it)
- [x] Cleaned up mocks: deleted `mockConnections.ts`, removed `MOCK_USERS` data (kept `mockUsers.ts` only for the `UserProfile` type)
- [x] Bigger/clearer Connections tabs (css)

---

## Pending (my work, not blocked)

- [ ] Side panel in Connections (like Discover's `ProfilePanel`)
      - Clicking a card should open the full-profile side panel, same as Discover
      - Needs `ProfilePanel` to be made variant-aware (its buttons are currently fixed to Discover's Connect/Dismiss). Make it swap buttons per context like UserCard does.
- [ ] Profile picture UPLOAD (Ivan added the endpoints)
      - `POST /me/picture` (JPEG/PNG up to 1MB, backend crops/resizes/cleans) -> returns new picture_url
      - `DELETE /me/picture` -> resets to default placeholder
      - Replace the current "paste a URL" field with a real file upload button
- [ ] `mockUsers.ts` is misnamed now (only holds the `UserProfile` type, no mocks). Rename to `types.ts` someday — low priority, touches many imports.

---

## Pending (waiting on backend or decisions)

- [ ] Auth: migrate to HttpOnly cookies (Ivan flagged this)
      - Backend stores JWT in an HttpOnly cookie now. Frontend should stop using localStorage + manual `Authorization: Bearer`, and use `credentials: 'include'` instead so the browser sends the cookie automatically. More secure (protects against XSS).
      - Touches `api.ts`, `connections.ts`, `profile.ts`, `auth.ts`, `updateMyLocation`, `respondToConnection`, etc. (everywhere we send the token by hand)
      - Ivan noted CORS config may need adjusting since front/back are on different origins. Not critical but recommended.
- [ ] Discover `match_score` still 0 (backend doesn't send a score yet; algorithm being designed with Ivan)
- [ ] Sports/modes from backend: `/activities` endpoint is active — can swap the hardcoded SPORTS list to fetch it
- [ ] Backend also blocking private endpoints if profile incomplete (security) — Ivan's call

---

## New work (build)

- [ ] Public pages navigation: Login/Register have no way back to landing except browser back button. Add logo->landing + link between login/register.
- [ ] Chat: needs WebSocket (real-time). Chat list, messages, unread badge, typing, online status. Ivan started `chat.go`. Not started on frontend.

---

## Match algorithm (design discussion with Ivan, ongoing)

- Hard filters already done in backend: location within radius + at least one shared sport.
- Score/ranking is the TODO. Agreed ideas: shared sports (weight), mode compatibility (weight, NOT filter — "Open to anything" is a wildcard), experience closeness (medium), distance-within-radius (matters more for some sports like running, less for gym/climbing), age by broad ranges (low).
- NOTE: interest_level is fixed at 3 for everyone, so it can't differentiate — don't use it in scoring.

---

## Decisions settled with Ivan

- [x] Dating mode + Gender: dropped (training-partners only)
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

- Recommendations active. `/activities` active. Connections fully done + tested. 100 seed users across Finland. Profile picture upload (MinIO S3) done. Chat started (`chat.go`).
- Auth: moving to cookie model, token expiry bumped for dev.

---

## Reminder

snake_case everywhere: `access_token`, `picture_url`, `activity_id`, `interaction_mode`, `max_radius`, `experience_level`, `interest_level`, `lat`, `lon`, `to_user_id`, `from_user_id`.