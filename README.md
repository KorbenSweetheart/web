<p align="center">
<img src="https://img.shields.io/badge/kood/Sisu-brightgreen?logo=gitea&logoColor=white&labelColor=8A2BE2">
<img src="https://img.shields.io/badge/ES2025-brightgreen?logo=typescript&logoColor=lemon&labelColor=white">
<img src="https://img.shields.io/badge/license-MIT-blue.svg">
</p>

# Match-me Web

Is a full-stack recommendation application, to connect users based on their profile information.

> [!NOTE]
> Project Context: Developed in 2026 as part of the JavaScript Module within a project-based coding program.

---

## Key Learnings & Skills Acquired

Through this task, we have learned and practically applied the following concepts:

<!-- TODO: update before submitting -->

- REST
- Full stack application
- React
- Typescript
- Uploading images
- Recommendations
- Realtime programming
- Security
- JWT
- Responsive web apps

## Project Scope & Development Constraints

<!-- TODO: update before submitting -->

These projects were built under specific course constraints to ensure mastery of the fundamentals:

- **Core Features (Mandatory):** Successfully implemented all essential requirements to ensure foundational functionality.
- **Enhancements (Extra):** Added optional features to improve the application and expand its capabilities.
- **Bonus Functionality:** Integrated innovative, out-of-scope features. These are managed via feature toggles to ensure the core functionality remains stable during reviews.
- **Library Constraints:** Relied primarily on standard JavaScript tools, using external libraries only when explicitly permitted by the project brief.

## Projects Included

- **Hello JS:** Start with NodeJS by creating `hello-world.js` to output "Hello, world!" into the console.
- **Real JS:** Build `ancient-history.js` to classify dates and define time.
- **Browser JS:** Use `get-el.js` to master element retrieval by tag, class, ID, or attribute, becoming an HTML navigator!

---

## Installation

### Prerequisites

<!-- TODO: update before submitting -->

- **Docker:** Version 29.6.1 or higher.
<!-- - **Go:** Version 1.24.0 or higher.
- **Node.js:** Version 20.6.0 or higher (required for native `--env-file` support in scripts; Node.js v22.x+ recommended for `--watch` stability).
- **npm:** Version 10.0.0 or higher. -->

## Quick Start

### 1. Clone the project

```bash
git clone https://gitea.kood.tech/ivanandreev/web
cd web
```

### 2. Create .env Files

Create a `.env` file inside the `backend/` and `frontend/` directories and configure the variables if you need:
For simplicity, you can rename the prepared-for-you file `.env.example` into `.env`.

```env
# Global variables
ENVIRONMENT=dev
JWT_SECRET=super-secure-random-key
...
```

### 3. Build the Dev Container

```bash
docker compose up --build
```

#### Docker useful commands

```bash
# Builds, (re)creates, starts, and attaches to containers for a service.
docker compose up

#Stops running containers without removing them.
docker compose stop

# Stops containers and removes containers, networks, volumes, and images created by up.
docker compose down
# Remove named volumes declared in the "volumes" section of the Compose file and anonymous volumes attached to containers
docker compose down -v

# Lists containers for a Compose project, with current status and exposed ports.
docker compose ps
```

### 4. Access the application

<!-- TODO: update before submitting -->

Open the following URL in your browser:

```
http://localhost:8080
```

[!INFO]

> To comply with the test case **"The interfaces are reachable by devices on other networks (not just localhost)."** We can launch an ngrok tunnel.
> Please contact the submitter when it is needed, since the app should be run on a local machine.

### 6. Testing and Examples files

<!-- TODO: update before submitting -->

We have prepared testing users for you that can use to speed up the testing process.

---

## User Guide

<!-- TODO: update before submitting -->

---

> [!CAUTION]
>
> ## ⚖️ Academic Integrity Disclaimer
>
> This project was submitted as part of the JavaScript Module within a project-based coding program. It is intended for portfolio purposes only.
>
> If you are a student currently enrolled in a similar course, please be aware that using any part of this code may violate your institution's academic honesty policy.
>
> Use this as a reference, but write your own code!

## License

This project is licensed under the [MIT License](https://gitea.kood.tech/ivanandreev/web/src/branch/main/LICENSE.txt).
