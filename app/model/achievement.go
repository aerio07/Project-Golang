package model

import "time"

type Achievement struct {
	ID        string    `json:"id"`
	StudentID string    `json:"student_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
