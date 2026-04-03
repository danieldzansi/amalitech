package main

import (
	"log"

	"github.com/danieldzansi/amalitech/handlers"
	"github.com/danieldzansi/amalitech/store"
	"github.com/gin-gonic/gin"
)

func main() {
	idempotencyStore := store.NewIdempotencyStore()

	r := gin.Default()

	r.POST("/payment/unsafe", handlers.WithoutIdempotency)
	r.POST("/payment/safe", handlers.WithIdempotency(idempotencyStore))

	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}