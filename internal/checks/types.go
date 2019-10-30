package checks

// Issue which is will be exposed.
type Issue struct {
	Namespace   string `json:"namespace,omitempty"`
	Name        string `json:"name"`
	Issue       string `json:"issue"`
	Description string `json:"description"`
	Command     string `json:"command"`
}
