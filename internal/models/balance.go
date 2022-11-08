package models

type Balance struct {
	UserId int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}
