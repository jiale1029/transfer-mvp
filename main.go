package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jiale1029/transaction/api"
	"github.com/jiale1029/transaction/dal"
	"github.com/jiale1029/transaction/dal/mysql"
)

const serverPort = ":8080"

func main() {
	initDAL()
	r := gin.Default()
	// accounts related endpoint
	r.GET("/accounts/:id", api.HandleGetAccount)
	r.POST("/accounts", api.HandleCreateAccount)
	r.GET("/accounts/list", api.HandleListAccounts)

	// transactions related endpoint
	r.GET("/transactions/:id", api.HandleGetTransaction)
	r.POST("/transactions", api.HandleSubmitTransaction)
	r.GET("/transactions/list", api.HandleListTransactions)
	r.Run(serverPort)
}

func initDAL() {
	api.AccountManager = dal.NewAccountDAO(mysql.InitMySQL())
}
