package lsp

type ShutdownRequest struct {
	Request
}

type ShutdownResponse struct {
	Response
	Result map[string]any `json:"result"`
}

type ExitNotification struct {
	Notification
}
