package models 

type PaymentRequest struct {
	IdempotencyKey string `json:"idempotency_key"`
	UserId string `json:"user_id"`
	Amount float64 `json:"amount"`
	currency string `json:"currency"`
}


type PaymentResponse struct{
	TransactionID string `json:"transaction_id"`
	UserId string `json"user_id"`
	Amount float64 `json:"amount"`
	Currency string `json:"currency"`
	Status string `json:"status"`
	Message string `json:"message"`
}