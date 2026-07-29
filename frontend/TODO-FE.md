1. Checklist de pendientes (para tu doc)

Esto es lo que queda por conectar/cambiar cuando Iván termine su parte:

Redirect después del login (la bandera de perfil):

LoginPage.tsx → ahora manda a todos a /app/profile. Cuando Iván dé la bandera profile_completed (o el endpoint /me/profile), cambiar para que: perfil completo → /app/discover, perfil vacío → /app/profile.

Guardar el perfil (profile.ts sigue con mock):

services/profile.ts → MOCK = true a false cuando Iván active /me/profile (PUT o POST). El fetch ya está escrito, solo se enciende.
Confirmar el método (PUT/POST) y que acepta la forma con activities: [{activity_id, experience, interest_level}].

Lista de deportes real:

services/profile.ts → reemplazar el array SPORTS fijo por un fetch a su endpoint de actividades, cuando Iván cargue los deportes en la BD.
Igual con MODES si los guarda en BD (tabla InteractionMode).

El token en las peticiones privadas:

Cuando conectes cualquier ruta privada (/me, /me/profile, recommendations...), cada fetch debe mandar Authorization: Bearer <token> leyendo el token de localStorage. Vale la pena crear services/api.ts que lo haga automático (para no repetirlo en cada llamada).

Logout de verdad:

Iván tiene /auth/logout comentado. Cuando lo active, tu logout() puede llamarlo además de borrar el token local.

Recommendations / Connections / Chat:

Trabajo nuevo (no cambios): services/users.ts que combine los 3 endpoints, UserCard, las páginas. Para cuando Iván tenga /recommendations, /users/:id, /users/:id/bio.

Recordatorio permanente: todo en snake_case (access_token, picture_url, activity_id, max_radius, profile_completed).