package lib

import "time"

type Site struct {
	ResourceID string                 `json:"id"`
	Name       string                 `json:"name"`
	CreatedAt  time.Time              `json:"created_at"`
	Meta       map[string]interface{} `json:"meta_data"`
}
