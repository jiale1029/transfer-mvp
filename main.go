package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jiale1029/transaction/api"
	"github.com/jiale1029/transaction/dal"
	"github.com/jiale1029/transaction/dal/mysql"
)

const serverPort = ":8080"

func main() {
	api.AccountManager = dal.NewAccountDAO(mysql.InitMySQL())
	r := gin.Default()
	// accounts related endpoint
	r.GET("/accounts/:id", api.HandleGetAccount)
	r.POST("/accounts", api.HandleCreateAccount)
	r.GET("/accounts/list", api.HandleListAccounts)

	// transactions related endpoint
	r.GET("/transactions/:id", api.HandleGetTransaction)
	r.POST("/transactions", api.HandleSubmitTransaction)
	r.GET("/transactions/list", api.HandleListTransactions)
	log.Fatal(r.Run(serverPort))
}
