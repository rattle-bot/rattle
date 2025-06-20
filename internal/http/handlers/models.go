package handlers

type Res struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type telegramInput struct {
	CheckDataString string `json:"check_data_string" validate:"required"`
	Hash            string `json:"hash" validate:"required"`
}

type createUserInput struct {
	TelegramID string `json:"telegram_id" validate:"required"`
	Role       string `json:"role" validate:"required,oneof=admin user"`
}

type createChatInput struct {
	ChatID string `json:"chat_id" validate:"required"`
	Send   bool   `json:"send"`
}

type updateChatInput struct {
	Send bool `json:"send"`
}

type saveContainerInput struct {
	Type  string `json:"type" validate:"required,oneof=name image id"`
	Value string `json:"value"`
	Mode  string `json:"mode" validate:"required,oneof=blacklist whitelist"`
}

type getRunningContainer struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Image   string            `json:"image"`
	Labels  map[string]string `json:"labels"`
	ShortID string            `json:"short_id"`
}

type createLogInput struct {
	Pattern   string `json:"pattern" validate:"required,min=1"`
	MatchType string `json:"match_type" validate:"required,oneof=include exclude"`
	EventType string `json:"event_type" validate:"required,oneof=error info warning success critical"`
}

type updateLogInput struct {
	Pattern   *string `json:"pattern" validate:"omitempty,min=1"`
	MatchType *string `json:"match_type" validate:"omitempty,oneof=include exclude"`
	EventType *string `json:"event_type" validate:"omitempty,oneof=error info warning success critical"`
}

type updateModeInput struct {
	Value string `json:"value" validate:"required,oneof=blacklist whitelist"`
}
