# TODO:

- [x] Add Router
- [x] Connect DB
  - [x] Add AutoMigrate
  - [x] Add seed data
  - [ ] Add 100-1000 seeded users
- [ ] Add Public Handlers
  - [x] Registration
    - [x] Add bcrypt (salt is the part of the algorithm)
    - [x] Add switch for different types of error, to return different statuses.
    - [ ] Add Login in or redirect to login upon successful registration
  - [x] Login
    - [x] Add JWT and session cookies
    - [ ] Add refresh token endpoint
  - [ ] Logout???
- [ ] Add Private Handlers
  - [x] Add /users/{id}
