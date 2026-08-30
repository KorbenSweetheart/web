# Test Cases

## Mandatory

- [x] It is possible to register with an email address and password.
- [x] The user can log out.
- [x] The application works with a single user.
- [x] It refuses to recommend an obviously poor match.
    Create two users in an empty system with obviously poor matching characteristics. Check to make sure that they are not recommended.
- [x] It recommends obviously good matches.
    Create two users in an empty system, who appear like they should obviously match.
- [x] The user is not shown any recommendations until they have completed their profile.
- [x] The user has a minimum of 5 biographical points to configure.
- [x] The user can change their biographical data.
- [x] The user can specify preference which target biographical data points.
- [x] A profile picture can be set.
- [x] The profile picture can be removed or changed.
- [x] The email address is not shown, except to the owner of the profile.
    The email address is not returned in API calls for other users.
- [x] The user can specify a location or preferred distance for matches.
- [x] The user only sees recommendations from their location.
- [x] The user can see a list of no more than 10 recommendations at a time.
- [x] The recommendations are prioritized with the best first.
- [x] The recommendations behave in line with the student's described matching logic.
- [x] It is possible to dismiss a recommendation.
    That recommendation is not shown again after it is dismissed.
- [x] Connection requests can be sent.
- [x] Incoming connection requests can be rejected.
- [x] Incoming connection requests can be accepted.
- [x] Users can only see profile information when properly allowed.
    They are recommended.
    There is an outstanding connection request.
    They are connected.
- [x] It is possible to disconnect with a user.
- [x] Chat is only possible between connected profiles.
- [x] Chats are ordered with the most recently active chat first.
- [x] Chat messages feature a date and time.
- [x] A chat history can be reached from the connected user's profile.
- [x] Both users see the same chat history.
- [x] The chat history API data is paginated.
- [x] The chat works in real time.
- [x] An unread message icon appears when new chat messages are received in real time.
- [x] The realtime implementation does not rely on polling.
- [ ] The recommendations endpoint only returns a list of ids. `Note: it returns score alongside IDs`
- [x] The connections endpoint only returns a list of ids.
- [x] The users endpoint returns a name and profile link.
- [x] The profile endpoint returns "about me" type information.
- [x] The bio endpoint returns biographical data.
- [x] All user responses return an id in the payload.
- [x] The me endpoints correctly shortcuts to the appropriate users endpoint.
- [x] The users endpoints return HTTP404 when the id is not found.
    This includes when the user is not allowed to see a profile. This is not quite how HTTP404 is described, but it means that a bad actor cannot distinguish between "does not exist", and "has blocked the user".
- [x] The backend is implemented using primary language from Coding Fundamentals.
- [x] The frontend is implemented in React using Typescript.
- [x] A PostgreSQL database is used as the primary application database.
- [x] The application is secure. Information is appropriately shown to the correct authenticated users only.
- [x] The application is responsive for mobile and desktop browsers.
- [x] A method was provided to load fictitious users into the system (minimum 100).

## Extra
- [x] The user experience is excellent, usable and well designed.
- [x] An offline/online indicator is shown on profile and chat views.
- [x] A typing in progress indicator is shown.
- [x] The recommendation algorithm is exceptional.
- [x] It implements proximity-based location filtering.