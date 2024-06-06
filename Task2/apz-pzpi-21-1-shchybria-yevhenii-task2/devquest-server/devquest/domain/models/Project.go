package models

import "github.com/google/uuid"

type Project struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Company Company `json:"company"`
}