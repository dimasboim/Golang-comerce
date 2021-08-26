package main

import (
	"Day15/config"
	"Day15/core"
	"Day15/midleware"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	config.InitDB()
	sqlDB, err := config.Db.DB()
	if err != nil {
		panic("failed db object")
	}
	defer sqlDB.Close()
	r := gin.Default()
	stock := r.Group("/stock")

	{
		stock.GET("/", midleware.IsAdmin(), core.GetStock)
		stock.GET("/:sku", midleware.IsAdmin(), core.GetStockBySKU)
		stock.POST("/insert", midleware.IsAdmin(), core.InsertProduct)
		stock.PUT("/restock", midleware.IsAdmin(), core.Restock)
		stock.PUT("/update", midleware.IsAdmin(), core.Update)
		stock.DELETE("/delete/:sku", midleware.IsAdmin(), core.Delete)
	}
	display := r.Group("/display")

	{
		display.GET("/", core.GetStockDisplay)
		display.GET("/:sku", core.GetStockBySKUDisplay)
		display.POST("/insert", midleware.IsAdmin(), core.InsertProductDisplay)
		display.PUT("/restock", midleware.IsAdmin(), core.RestockDisplay)
		display.PUT("/update", midleware.IsAdmin(), core.UpdateDisplay)
		display.DELETE("/delete/:sku", midleware.IsAdmin(), core.DeleteDisplay)
	}
	 
	auth := r.Group("/auth")
	{
		 	auth.POST("/register", core.Register)
		auth.POST("/login", core.Login)
		auth.POST("/logout",midleware.IsUser(), core.Logout)
		auth.POST("/logoutAdmin",midleware.IsAdmin(), core.Logout)
		 
	}

	cart := r.Group("/cart")
	{
		cart.GET("/", midleware.IsUser(), core.Getcart)
		cart.POST("/add", midleware.IsUser(), core.Addtocart)
		cart.DELETE("/delete/:sku", midleware.IsUser(), core.DeleteBySKU)

		cart.POST("/checkout", midleware.IsUser(), core.CheckoutCart)
	}
	trans := r.Group("/transaction")
	{
		trans.GET("/", midleware.IsUser(), core.GetTransactionHistory)

	}
	report := r.Group("/report")
	{
		report.GET("/", midleware.IsAdmin(), core.GetAllReport)

	}
	r.Run(":" + os.Getenv("API_PORT"))
}
