package src

type IncomeMessage struct {
	Command string `json:"command"`
	Room    string `json:"room"`
	Message string `json:"message"`
}

type SendMessage struct {
	Room    string `json:"room"`
	Message string `json:"message"`
}
