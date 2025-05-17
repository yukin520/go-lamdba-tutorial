package domain

import "time"

type ToDo struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CraetedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
