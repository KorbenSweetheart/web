# REST API for Match Me Application

## Technical Requirements:

- Implemented in Go or Typescript.
- Data must be persisted in a PostgreSQL database.
- The application must be secure. It must not leak private information, or allow access to data which must not be seen.
- Passwords should be encripted with bcrypt & salt
- Sessions are managed with JWT

## RESTful endpoints.

> [!CAUTION]
> None of them must return authentication-related data.

- `/users/{id}`: which returns the user's name and link to the profile picture.
- `/users/{id}/profile`: which returns the users "about me" type information.
- `/users/{id}/bio`: which returns the users biographical data (the data used to power recommendations).
- `/me`: which is a shortcut to `/users/{id}` for the authenticated user. You should also implement `/me/profile` and `/me/bio`.
- `/recommendations`: which returns a maximum of 10 recommendations, containing only the `id` and nothing else.
- `/connections`: which returns a list connected profiles, containing only the `id` and nothing else.

All of the responses for `/users` data must also contain the `id`.

If the `id` is not found, or the user does not have permission to view that profile, it must return `HTTP404`.

### Example usage

Let's say you wanted to show a list of recommendations with icons reflecting their bio, then you'd need to:

1. fetch `/recommendations` to get a list of ids
2. fetch `/users/{id}` for name and profile picture for each user
3. fetch `/users/{id}/bio` for biographical information for each user
4. Combine all the data into one object.

## Synthetic Data for Review

Load fictitious users into the system, with different profiles, to see how matching works with some scale. A minimum of 100 users.

It must be possible to drop database, and reload those users as a separate step so the reviewer can see your application working with no users, few users or many users.
