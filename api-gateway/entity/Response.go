package entity

type Response struct {
	Status  int         `json:"status,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
