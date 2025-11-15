package model

type ApiResponse struct {
    Message string `json:"message"`
    Status  int    `json:"status"`
    Success bool   `json:"success"`
    Data    any    `json:"data,omitempty"`
}