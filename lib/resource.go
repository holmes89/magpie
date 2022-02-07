package lib

import "time"

type Resource struct {
	Type       string                 `json:"type"`
	ResourceID string                 `json:"id"`
	Name       string                 `json:"name"`
	CreatedAt  time.Time              `json:"created_at"`
	Meta       map[string]interface{} `json:"meta_data"`
}