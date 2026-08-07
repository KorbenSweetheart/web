package handlers

// Handlers combine prepared HTTP handlers to pass them into the router
type Handlers struct {
	Auth  *AuthHandler
	User  *UserHandler
	Match *RecommendationsHandler
}
