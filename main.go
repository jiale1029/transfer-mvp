package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jiale1029/transaction/api"
)

func main() {
	r := gin.Default()

	// accounts related endpoint
	r.GET("/accounts/:id", api.HandleGetAccount)
	r.POST("/accounts", api.HandleCreateAccount)
	r.GET("/accounts/list", api.HandleListAccounts)

	// transactions related endpoint
	r.GET("/transactions/:id", api.HandleGetTransaction)
	r.POST("/transactions", api.HandleSubmitTransaction)
	r.GET("/transactions/list", api.HandleListTransactions)

	log.Fatal(r.Run())
}
