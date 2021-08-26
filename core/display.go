package core

import (
	"Day15/config"
	"Day15/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStockDisplay(c *gin.Context) {
	stockData := []models.Product_display{}
	config.Db.Find(&stockData)
	c.JSON(200, gin.H{
		"status": "Success",
		"data":   stockData,
	})
}
func GetStockBySKUDisplay(c *gin.Context) {
	stockData := []models.Product_display{}
	sku := c.Param("sku")
	config.Db.First(&stockData, "sku = ?", sku)
	//	Db.Find(&blogData)
	c.JSON(200, gin.H{
		"status": "Success",
		"data":   stockData,
	})
}
func InsertProductDisplay(c *gin.Context) {
	if c.PostForm("sku") == "" || c.PostForm("price") == "" || c.PostForm("qty") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Missing required parameter",
		})
		return
	}

	prc, err := strconv.ParseFloat(c.PostForm("price"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid Format Price",
		})
	}

	qt, err := strconv.ParseInt(c.PostForm("qty"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid Format Qty",
		})
	}

	var (
		sku              = c.PostForm("sku")
		stockData        models.Product_warehouse
		stockDataDisplay models.Product_display
	)

	errz := config.Db.First(&stockData, "sku = ?", sku)
	if errz.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found sku in warehouse",
			"Message": errz.Error.Error(),
			"data":    stockData,
		})
		return
	}
	if stockData.Qty < qt {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Refuse, qty > warehouse",
			"Message": errz.Error.Error(),
			"data":    stockData,
		})
		return
	}
	if int64(stockData.Price) > int64(prc) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Refuse, price display < warehouse",
			"Message": errz.Error.Error(),
			"data":    stockData,
		})
		return
	}
	stockData.Qty = stockData.Qty - qt

	ex := config.Db.Save(&stockData)
	if ex.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"Message": ex.Error.Error(),
			"data":    stockData,
		})
		return
	}
	name := stockData.Name
	stockDataDisplay.Name = name
	stockDataDisplay.Qty = qt
	stockDataDisplay.Sku = sku
	stockDataDisplay.Price = prc
	stockDataDisplay.User_id = c.GetUint("user_id")

	er := config.Db.Create(&stockDataDisplay)
	if er.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"Message": er.Error.Error(),
			"data":    stockDataDisplay,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Ok",
			"Message": "Data Insert Stock Display",
			"data":    stockDataDisplay,
		})

	}
}

func RestockDisplay(c *gin.Context) {
	if c.Query("sku") == "" || c.Query("qty") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Missing required parameter",
		})
		return
	}

	qt, err := strconv.ParseInt(c.Query("qty"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid Format Qty",
		})
	}

	var (
		sku              = c.Query("sku")
		qty              = qt
		stockData        models.Product_warehouse
		stockDataDisplay models.Product_display
	)
	errz := config.Db.First(&stockData, "sku = ?", sku)
	if errz.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found sku in warehouse",
			"Message": errz.Error.Error(),
			"data":    stockData,
		})
		return
	}
	if stockData.Qty < qt {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Refuse, qty > warehouse",
			"Message": errz.Error.Error(),
			"data":    stockData,
		})
		return
	}

	stockData.Qty = stockData.Qty - qt

	errw := config.Db.First(&stockDataDisplay, "sku = ?", sku)
	if errw.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"Message": errw.Error.Error(),
			"data":    stockDataDisplay,
		})
		return
	}
	stockDataDisplay.Qty = stockDataDisplay.Qty + qty

	errx := config.Db.Save(&stockDataDisplay)
	if errx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"Message": errx.Error.Error(),
			"data":    stockDataDisplay,
		})
		return
	}
	ex := config.Db.Save(&stockData)
	if ex.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"Message": ex.Error.Error(),
			"data":    stockData,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Ok",
		"Message": "Data Updated Restock",
		"data":    stockDataDisplay,
	})

}
func UpdateDisplay(c *gin.Context) {
	if c.Query("sku") == "" || c.Query("name") == "" || c.Query("price") == "" || c.Query("qty") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Missing required parameter",
		})
		return
	}

	qt, err := strconv.ParseInt(c.Query("qty"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid Format Qty",
		})
	}
	prc, err := strconv.ParseFloat(c.Query("price"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid Format Price",
		})
	}
	var (
		sku       = c.Query("sku")
		qty       = qt
		name      = c.Query("name")
		price     = prc
		stockData models.Product_display
	)
	errz := config.Db.First(&stockData, "sku = ?", sku)
	if errz.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"Message": errz.Error.Error(),
			"data":    stockData,
		})
		return
	}
	stockData.Qty = qty
	stockData.Price = price
	stockData.Name = name

	errx := config.Db.Save(&stockData)
	if errx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"Message": errx.Error.Error(),
			"data":    stockData,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Ok",
		"Message": "Data Updated ",
		"data":    stockData,
	})

}
func DeleteDisplay(c *gin.Context) {
	if c.Param("sku") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Missing required parameter",
		})
		return
	}

	var (
		sku              = c.Param("sku")
		stockData        models.Product_warehouse
		stockDataDisplay models.Product_display
	)
	err := config.Db.First(&stockData, "sku = ?", sku)
	if err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"Message": err.Error.Error(),
			"data":    stockData,
		})
		return
	}
	err = config.Db.First(&stockDataDisplay, "sku = ?", sku)
	if err.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"Message": err.Error.Error(),
			"data":    stockData,
		})
		return
	}
	stockData.Qty = stockData.Qty + stockDataDisplay.Qty
	err = config.Db.Delete(&stockDataDisplay)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"Message": err.Error.Error(),
			"data":    stockDataDisplay,
		})
		return
	}

	err = config.Db.Save(&stockData)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"Message": err.Error.Error(),
			"data":    stockDataDisplay,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Ok",
		"Message": "Data deleted display",
		"data":    stockData,
	})

}
