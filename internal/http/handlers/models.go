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
	Send   *bool  `json:"send"`
}

type updateChatInput struct {
	Send bool `json:"send" validate:"required"`
}
