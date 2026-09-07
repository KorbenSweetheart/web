<p align="center">
<img src="https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white">
<img src="https://img.shields.io/badge/Echo-v5-00ADD8">
<img src="https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black">
<img src="https://img.shields.io/badge/TypeScript-5.x-3178C6?logo=typescript&logoColor=white">
<img src="https://img.shields.io/badge/PostgreSQL-18-4169E1?logo=postgresql&logoColor=white">
<img src="https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white">
</p>

# Pulse

Pulse is a full-stack recommendation and real-time social platform that helps sports enthusiasts find compatible training partners based on their athletic profiles, shared activities, experience levels, and geographical proximity.

> [!NOTE]
> **Overfetching is a mandatory requirement in this project. **
> 
> An additional task is to implement GraphQL to display its value in contrast.

> [!NOTE]
> Project Context: Developed in 2026 as part of the JavaScript Module within a project-based coding program.

---

## Key Learnings & Skills Acquired

Through this project, the following technical concepts and architectural patterns were applied:

- Full-stack system architecture with separation of concerns between Go backend and React frontend.
- RESTful API design with strict permission isolation and secure 404 responses for unauthorized resource access.
- High-performance web service development using Go and the Echo v5 web framework.
- Relational data modeling, relationships, and queries using PostgreSQL and GORM ORM.
- Proximity-based recommendation filtering using GPS coordinates and configurable search radii.
- Real-time bidirectional communication using WebSockets for live chat, presence tracking, and typing indicators.
- S3-compatible object storage integration using MinIO for profile image upload and retrieval.
- Secure authentication architecture with bcrypt password hashing and dual-token JWT management.
- Modern frontend single-page application built with React 19, TypeScript, and Vite.
- Responsive layout design and UI state management without relying on third-party UI component frameworks.

---

## Authentication & Session Management

The platform implements a secure dual-token authentication architecture:

- **Password Security:** Passwords are never stored in plain text; they are hashed using bcrypt with a salt and a cost factor of 12.
- **Access Tokens:** Short-lived JSON Web Tokens (JWT) signed with HMAC-SHA256 containing custom claims (`user_id`). They are delivered via HttpOnly cookies (`access_token`) and supported through the `Authorization: Bearer <token>` header for API clients.
- **Refresh Tokens:** Long-lived, cryptographically secure 32-byte random tokens generated using crypto/rand. To protect against database compromise, only the SHA-256 hash of the token is persisted in PostgreSQL. Delivered via an HttpOnly cookie scoped specifically to `/auth/refresh`.
- **Token Rotation:** Every call to the `/auth/refresh` endpoint validates the token hash and expiration in the database, invalidates the used refresh token, and issues a new access token and refresh token pair.
- **Logout & Invalidation:** Logging out deletes the stored refresh token from the database by account ID and clears all authentication cookies from the client browser.

---

## Project Scope & Development Constraints

The application was built under structured course specifications to ensure mastery of foundational technologies:

- **Core Features (Mandatory):**
  - Registration with unique email validation and bcrypt password hashing.
  - Login, logout, and session lifecycle managed by JWT.
  - Profile completion workflow: users cannot view recommendations or connect until their profile is complete.
  - Rich biographical profiles capturing at least 5 data points: name, age, bio, activity selections, experience levels, interest levels, interaction modes, and location.
  - Profile image upload, change, and removal via MinIO object storage.
  - Privacy safeguards: email addresses are strictly private (visible only to the owner on `/me`), and unauthorized profile queries return HTTP 404.
  - Recommendation engine prioritizing the strongest matches up to a maximum of 10 users at a time.
  - Dismissal persistence: dismissed recommendations are never shown again to the user.
  - Connection lifecycle: send request, accept request, reject request, and disconnect.
  - Real-time chat accessible between connected profiles with paginated history and message timestamps.
  - Unopinionated generic REST API with endpoints (`/users/{id}`, `/users/{id}/profile`, `/users/{id}/bio`, `/me`, `/recommendations`, `/connections`).

- **Enhancements (Extra Requirements):**
  - **Online/Offline Status:** Real-time presence indicator displayed on user profile pages and active chat views.
  - **Typing Indicator:** Real-time typing status notification in active chat views that automatically clears when typing ceases.
  - **Unread Notification & Sorting:** Live unread message notification badge and dynamic reordering of chat conversations by most recent message.
  - **Proximity-Based Filtering:** Proximity matching utilizing browser geolocation coordinates, user-defined maximum search radius, and distance calculation.

- **Bonus Functionality:**
  - **S3 Object Storage Pipeline:** MinIO object storage container integration for media handling with read-only public access policies and default image fallbacks.
  - **Synthetic Data Generator:** Automated database migration and seeder capable of generating 1,000+ realistic synthetic profiles and pre-configured test scenarios.

---

## Architecture & Components

The application consists of modular services coordinated through containerization:

- **Backend (`backend/`):** Go service built on the Echo v5 framework, exposing REST endpoints and WebSocket handlers for real-time events.
- **Frontend (`frontend/`):** React 19 single-page application written in TypeScript, bundled with Vite, styled with custom CSS, and routed with React Router v7.
- **Database (`db`):** PostgreSQL database persisting accounts, user profiles, activity mappings, connection records, refresh token hashes, and chat messages.
- **Object Storage (`minio`):** MinIO S3-compatible service storing uploaded profile images.
- **Migrator (`migrator`):** Database initialization service that executes schema migrations, loads system activity dictionaries, and handles synthetic user seeding.

---

## Installation

### Prerequisites

- **Docker:** Version 20.10 or higher.
- **Docker Compose:** Version 2.0 or higher.

---

## Quick Start

### 1. Clone the project

```bash
git clone https://gitea.kood.tech/ivanandreev/web
cd web
```

### 2. Configure Environment Files

Create `.env` files in both the `backend/` and `frontend/` directories. You can copy the provided example templates:

```bash
cp backend/.env.example backend/.env
```

Key backend configuration options in `backend/.env`:

```env
ENVIRONMENT=dev
JWT_SECRET=super-secure-random-key
SERVER_ADDRESS=0.0.0.0:8080
DB_NAME=postgres
DB_USER=dbuser
DB_PASS=my_db_password
DB_HOST=db
DB_PORT=5432
MINIO_ROOT_USER=minio_admin
MINIO_ROOT_PASSWORD=minio_password
MINIO_ENDPOINT=minio:9000
MINIO_PUBLIC_ENDPOINT=http://localhost:9000
SEED_USERS=true
```

### 3. Build and Run the Application

Start all services using Docker Compose:

```bash
docker compose up --build
```

#### Useful Docker Commands

```bash
# Start all containers in the background
docker compose up -d

# Stop running containers without removing volumes
docker compose stop

# Stop containers and remove networks and containers
docker compose down

# Stop containers and wipe database and storage volumes (resets all data)
docker compose down -v

# View container statuses and port mappings
docker compose ps

# View backend or frontend logs
docker compose logs -f backend
docker compose logs -f frontend
```

### 4. Access the Application

Once the containers are healthy and running, access the services:

- **Frontend Application:** `http://localhost:5173`
- **Backend REST API:** `http://localhost:8080`
- **MinIO Web Console:** `http://localhost:9001` (Username: `minio_admin`, Password: `minio_password`)
- **PostgreSQL Database:** `localhost:5432`

---

## Reviewer Testing & Synthetic Data Guide

### Testing Edge Cases with an Empty Database

To test edge cases such as single-user behavior, zero-recommendation states, and matching validation with isolated user pairs:

1. Set `SEED_USERS=false` in `backend/.env` (or in `docker-compose.yml`).
2. Wipe existing volumes and start the containers:
   ```bash
   docker compose down -v && docker compose up --build
   ```
3. The database will initialize with dictionaries only and zero accounts.
4. **Poor Match Edge Case:** Register User A and User B with conflicting interaction preferences, completely disjoint activities, and distant locations. Verify that neither user appears in the other's recommendations list (`/recommendations`).
5. **Good Match Edge Case:** Register User C and User D with identical or overlapping activities, compatible experience levels, and within each other's search radius. Verify that they receive high recommendation scores and appear in each other's discovery feed.

### Testing at Scale with Seeded Accounts

When `SEED_USERS=true`, the migrator seeds over 1,000 synthetic profiles across multiple cities along with pre-configured test accounts.

All seeded test accounts share the default password: `12345678`.

| Email | Name | Focus / Activities | Location | Pre-configured State |
| :--- | :--- | :--- | :--- | :--- |
| `obiwan@matchme.com` | Obi-Wan Kenobi | Running, Aikido, Yoga | Helsinki | Connected to Anakin, Yoda; Pending from Jar Jar; Declined Maul |
| `anakin@matchme.com` | Anakin Skywalker | CrossFit, MMA, Running | Espoo | Connected to Obi-Wan, Ahsoka; Declined Dooku |
| `yoda@matchme.com` | Master Yoda | Yoga, Aikido | Helsinki | Connected to Obi-Wan, Windu; Pending from Ahsoka |
| `windu@matchme.com` | Mace Windu | MMA, Gym | Helsinki | Connected to Yoda |
| `dooku@matchme.com` | Count Dooku | Aikido, Padel | Espoo | Connected to Ventress; Pending from Maul |
| `maul@matchme.com` | Darth Maul | Jiu-Jitsu, CrossFit, Climbing | Vantaa | Pending to Dooku; Declined by Obi-Wan |
| `ventress@matchme.com` | Asajj Ventress | MMA, Climbing | Helsinki | Connected to Dooku |
| `ahsoka@matchme.com` | Ahsoka Tano | Running, Cycling, Jiu-Jitsu | Helsinki | Connected to Anakin; Pending to Yoda |
| `jarjar@matchme.com` | Jar Jar Binks | Swimming, Running, Football | Helsinki | Pending to Obi-Wan |

---

## User Guide

### 1. Registration and Profile Completion

1. Navigate to `http://localhost:5173` and click **Register**.
2. Enter your name, email address, and a password (minimum 8 characters).
3. Upon first login, you are directed to the **Profile Setup** page. Recommendations and connection features remain locked until this setup is finished.
4. Fill in biographical details:
   - Bio and summary.
   - Age and interaction mode (Social, Focused, Open to anything).
   - Select sports and activities, setting your experience level and interest level for each.
   - Set location coordinates using browser geolocation or manual selection, and specify a maximum matching radius in kilometers.
   - Upload a profile picture (stored in MinIO) or proceed with the default avatar.
5. Save your profile to unlock all application features.

### 2. Discovering Matches

1. Navigate to the **Discover** page.
2. The recommendation engine evaluates user profiles within your location radius and calculates a match score based on shared activities, skill compatibility, and interaction mode.
3. The top 10 ranked recommendations are presented with their compatibility scores and profile summaries.
4. For each recommendation:
   - Click **Connect** to send a connection request.
   - Click **Dismiss** to remove the recommendation. Dismissed profiles will not be shown again.

### 3. Managing Connections

1. Navigate to the **Connections** page.
2. Review incoming connection requests: click **Accept** to establish a mutual connection or **Reject** to decline.
3. View your active connections list.
4. You can disconnect with any user at any time using the **Disconnect** option.

### 4. Real-Time Chat

1. Open a conversation directly from a connected user's profile or from the **Chats** page.
2. Messages are delivered in real time over WebSockets without polling.
3. Chat features:
   - **Presence:** View online and offline status in real time.
   - **Typing Indicator:** See when the other user is typing (`... typing`).
   - **Message History:** Message history is loaded in paginated batches with clear timestamps.
   - **Unread Notifications:** Unread indicators highlight chats with new messages, and conversations reorder automatically by the most recent message.

---

> [!CAUTION]
>
> ## Academic Integrity Disclaimer
>
> This project was submitted as part of the JavaScript Module within a project-based coding program. It is intended for portfolio purposes only.
>
> If you are a student currently enrolled in a similar course, please be aware that using any part of this code may violate your institution's academic honesty policy.
>
> Use this as a reference, but write your own code!

## License

This project is licensed under the [MIT License](https://gitea.kood.tech/ivanandreev/web/src/branch/main/LICENSE.txt).
