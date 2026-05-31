package plugin

type CommandRequest struct {
	RequestID string `json:"request_id"`
}

type CommandResponse struct {
	RequestID string `json:"request_id"`
	PluginID  string `json:"plugin_id"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}
