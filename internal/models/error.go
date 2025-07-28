package models

// Error codes
const (
	ErrCodeInvalidRequest = iota
	ErrCodeInvalidID
	ErrCodeNotFound
	ErrCodeDatabaseOperation
	ErrCodeInvalidDate
)

type Error struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Details string `json:"details"`
}
