package handlers

type Res struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type telegramDataInput struct {
	CheckDataString string `json:"check_data_string" validate:"required"`
	Hash            string `json:"hash" validate:"required"`
}
