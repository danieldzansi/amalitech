package handlers

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/danieldzansi/amalitech/models"
	"github.com/gin-gonic/gin"
)

var chargeCount int64

func WithoutIdempotency(c *gin.Context) {
	var req models.PaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	count := atomic.AddInt64(&chargeCount, 1)
	txID := fmt.Sprintf("TXN-UNSAFE-%d", count)

	fmt.Printf(
		"[NO IDEMPOTENCY] Charging user %s GHS %.2f | Transaction: %s\n",
		req.UserID, req.Amount, txID,
	)

	c.JSON(http.StatusOK, models.PaymentResponse{
		TransactionID: txID,
		UserID:        req.UserID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Status:        "success",
		Message:       "Payment processed (NO idempotency — duplicate charges possible!)",
	})
}

func GetChargeCount() int64 {
	return atomic.LoadInt64(&chargeCount)
}

func ResetChargeCount() {
	atomic.StoreInt64(&chargeCount, 0)
}