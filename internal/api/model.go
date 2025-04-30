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
