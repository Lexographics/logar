package api

type ArgType struct {
	Type string `json:"type"`
	Kind string `json:"kind"`
}

type ActionDetails struct {
	Path        string    `json:"path"`
	Args        []ArgType `json:"args"`
	Description string    `json:"description"`
}

type SessionData struct {
	Device       string `json:"device"`
	LastActivity string `json:"last_activity"`
	CreatedAt    string `json:"created_at"`
	IsCurrent    bool   `json:"is_current"`
	Token        string `json:"token"`
}
