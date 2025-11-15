package model

import "time"

type Item struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Done      bool      `json:"done"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type ListItemsResponse struct {
    Items []Item `json:"items"`
    Count int    `json:"count"`
}

type CreateItemRequest struct {
    Name string `json:"name"`
}

type UpdateItemRequest struct {
    Name *string `json:"name,omitempty"`
    Done *bool   `json:"done,omitempty"`
}