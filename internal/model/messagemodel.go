package model

import "time"

type MessageResponse struct {
	ID        uint      `json:"id"`
    From      uint      `json:"from"`
    To        uint      `json:"to"` 
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    IsRead    bool      `json:"is_read"`
}