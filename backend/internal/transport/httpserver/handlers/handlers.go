package handlers

// Handlers combine prepared HTTP handlers to pass them into the router
type Handlers struct {
	Health     *HealthHandler
	Dictionary *DictionaryHandler
	Auth       *AuthHandler
	User       *UserHandler
	Match      *MatchHandler
	Conn       *ConnectionHandler
}
