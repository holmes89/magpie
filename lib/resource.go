package lib

import (
	"time"
)

type Site struct {
	Type      string                 `json:"-" dynamodbav:"ID"`
	ID        string                 `json:"id" dynamodbav:"SK"`
	Name      string                 `json:"name" dynamodbav:"name"`
	URL       string                 `json:"url" dynamodbav:"url"`
	CreatedAt time.Time              `json:"created_at" dynamodbav:"created_at"`
	Meta      map[string]interface{} `json:"meta_data" dynamodbav:"meta_data"`
}
