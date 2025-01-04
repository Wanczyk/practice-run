package src

type StatusType string

const (
	StatusSuccess StatusType = "Success"
	StatusError   StatusType = "Error"
)

type Commands string

const (
	CreateCommand  Commands = "create"
	JoinCommand    Commands = "join"
	LeaveCommand   Commands = "leave"
	SendCommand    Commands = "send"
	UnknownCommand Commands = "unknown"
)

type IncomeMessage struct {
	Command Commands `json:"command"`
	Data    *Data    `json:"data"`
}

type SendMessage struct {
	Status  StatusType `json:"status"`
	Command Commands   `json:"command"`
	Data    *Data      `json:"data,omitempty"`
	Error   *Error     `json:"error,omitempty"`
}

type Data struct {
	Room    string `json:"room"`
	Message string `json:"message,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
