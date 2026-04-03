package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/danieldzansi/amalitech/handlers"
	"github.com/danieldzansi/amalitech/models"
	"github.com/danieldzansi/amalitech/store"
	"github.com/gin-gonic/gin"
)

const serverAddr = "http://localhost:8080"

func main() {
	gin.SetMode(gin.ReleaseMode)

	idempotencyStore := store.NewIdempotencyStore()
	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/payment/unsafe", handlers.WithoutIdempotency)
	r.POST("/payment/safe", handlers.WithIdempotency(idempotencyStore))

	go func() {
		if err := r.Run(":8080"); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	time.Sleep(300 * time.Millisecond)

	fmt.Println(strings.Repeat("=", 65))
	fmt.Println("       IDEMPOTENCY DEMO — Gin Payment API")
	fmt.Println(strings.Repeat("=", 65))

	runWithoutIdempotencyDemo()
	fmt.Println()
	runWithIdempotencyDemo()

	fmt.Println()
	fmt.Println(strings.Repeat("=", 65))
	fmt.Println("                     SUMMARY")
	fmt.Println(strings.Repeat("=", 65))
	fmt.Printf("Without Idempotency — Total charges made: %d  (overcharged!)\n", handlers.GetChargeCount())
	fmt.Printf("With Idempotency    — Total charges made: %d  (correct!)\n", handlers.GetSafeChargeCount())
	fmt.Println(strings.Repeat("=", 65))
}

func runWithoutIdempotencyDemo() {
	fmt.Println()
	fmt.Println("SCENARIO 1: POST /payment/unsafe  (WITHOUT IDEMPOTENCY)")
	fmt.Println("Client retries the same payment 3 times due to network lag")
	fmt.Println(strings.Repeat("-", 65))

	payment := models.PaymentRequest{
		UserID:   "user_daniel",
		Amount:   500.00,
		Currency: "GHS",
	}

	for i := 1; i <= 3; i++ {
		fmt.Printf("\n  [Attempt %d] Sending POST /payment/unsafe...\n", i)
		resp, err := postRequest("/payment/unsafe", payment)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}
		fmt.Printf("  Response → Status: %s | TxID: %s\n", resp.Status, resp.TransactionID)
	}

	fmt.Println()
	fmt.Printf("  💸 User was charged GHS 500 a total of %d times!\n", handlers.GetChargeCount())
}

func runWithIdempotencyDemo() {
	fmt.Println()
	fmt.Println("SCENARIO 2: POST /payment/safe  (WITH IDEMPOTENCY)")
	fmt.Println("    Client retries same payment 3 times — with an Idempotency-Key")
	fmt.Println(strings.Repeat("-", 65))

	payment := models.PaymentRequest{
		IdempotencyKey: "idem-key-abc123",
		UserID:         "user_daniel",
		Amount:         500.00,
		Currency:       "GHS",
	}

	for i := 1; i <= 3; i++ {
		fmt.Printf("\n  [Attempt %d] Sending POST /payment/safe (key: %s)...\n", i, payment.IdempotencyKey)
		resp, err := postRequest("/payment/safe", payment)
		if err != nil {
			fmt.Printf("  Error: %v\n", err)
			continue
		}
		fmt.Printf("  Response → Status: %s | TxID: %s\n", resp.Status, resp.TransactionID)
	}

	fmt.Println()
	fmt.Printf("User was charged GHS 500 exactly %d time — no matter how many retries!\n", handlers.GetSafeChargeCount())
}

func postRequest(path string, payload models.PaymentRequest) (models.PaymentResponse, error) {
	body, _ := json.Marshal(payload)

	res, err := http.Post(serverAddr+path, "application/json", bytes.NewReader(body))
	if err != nil {
		return models.PaymentResponse{}, err
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)

	var resp models.PaymentResponse
	json.Unmarshal(raw, &resp)
	return resp, nil
}