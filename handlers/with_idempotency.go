package handlers

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/danieldzansi/amalitech/models"
	"github.com/danieldzansi/amalitech/store"
	"github.com/gin-gonic/gin"
)

var safeChargeCount int64

func WithIdempotency(s *store.IdempotencyStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.PaymentRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if req.IdempotencyKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Idempotency-Key is required"})
			return
		}

		if cached, found := s.Get(req.IdempotencyKey); found {
			fmt.Printf(
				"[WITH IDEMPOTENCY] Duplicate key detected: %s | Returning cached result — user NOT charged again\n",
				req.IdempotencyKey,
			)
			c.Header("X-Idempotent-Replayed", "true")
			c.JSON(http.StatusOK, cached)
			return
		}
		count := atomic.AddInt64(&safeChargeCount, 1)
		txID := fmt.Sprintf("TXN-SAFE-%d", count)

		fmt.Printf(
			" [WITH IDEMPOTENCY] First-time request. Charging user %s GHS %.2f | Transaction: %s\n",
			req.UserID, req.Amount, txID,
		)

		resp := &models.PaymentResponse{
			TransactionID: txID,
			UserID:        req.UserID,
			Amount:        req.Amount,
			Currency:      req.Currency,
			Status:        "success",
			Message:       "Payment processed successfully",
		}

		s.Set(req.IdempotencyKey, resp)

		c.JSON(http.StatusOK, resp)
	}
}

func GetSafeChargeCount() int64 {
	return atomic.LoadInt64(&safeChargeCount)
}

func ResetSafeChargeCount() {
	atomic.StoreInt64(&safeChargeCount, 0)
}