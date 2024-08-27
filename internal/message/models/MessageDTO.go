package models

type MessageDTO struct {
	SenderID string `json:"sender_id"`
	Message  string `json:"message"`
}
