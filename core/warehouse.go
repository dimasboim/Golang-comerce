package core

import (
	"Day15/config"
	"Day15/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStock(c *gin.Context) {
	stockData := []models.Product_warehouse{}
	config.Db.Find(&stockData)
	c.JSON(200, gin.H{
		"status": "Success",
		"data":   stockData,
	})
}
func GetStockBySKU(c *gin.Context) {
	stockData := []models.Product_warehouse{}
	sku := c.Param("sku")
	config.Db.First(&stockData, "sku = ?", sku)
	//	Db.Find(&blogData)
	c.JSON(200, gin.H{
		"status": "Success",
		"data":   stockData,
	})
}
func InsertProduct(c *gin.Context) {
	if c.PostForm("sku") == "" || c.PostForm("name") == "" || c.PostForm("price") == "" || c.PostForm("qty") == "" {
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
		return
	}

	qt, err := strconv.ParseInt(c.PostForm("qty"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid Format Qty",
		})
		return
	}
	stockData := models.Product_warehouse{
		Sku:     c.PostForm("sku"),
		Name:    c.PostForm("name"),
		Price:   prc,
		Qty:     qt,
		User_id: c.GetUint("user_id"),
	}
	var stockDataexist models.Product_warehouse
	config.Db.First(&stockDataexist, "sku = ?", stockData.Sku)
	if stockData.Sku == stockDataexist.Sku {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Sku already exist",
		})
		return
	}

	errz := config.Db.Create(&stockData)
	if errz.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"Message": errz.Error.Error(),
			"data":    stockData,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Ok",
			"Message": "Data Insert Stock",
			"data":    stockData,
		})

	}
}

func Restock(c *gin.Context) {
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
		return
	}

	var (
		sku       = c.Query("sku")
		qty       = qt
		stockData models.Product_warehouse
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
		"Message": "Data Updated Restock",
		"data":    stockData,
	})

}
func Update(c *gin.Context) {
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
		return
	}
	prc, err := strconv.ParseFloat(c.Query("price"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid Format Price",
		})
		return
	}
	var (
		sku       = c.Query("sku")
		qty       = qt
		name      = c.Query("name")
		price     = prc
		stockData models.Product_warehouse
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
func Delete(c *gin.Context) {
	if c.Param("sku") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Missing required parameter",
		})
		return
	}

	var (
		sku = c.Param("sku")

		stockData models.Product_warehouse
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

	err = config.Db.Delete(&stockData)
	if err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"Message": err.Error.Error(),
			"data":    stockData,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Ok",
		"Message": "Data deleted Stock warehouse",
		"data":    stockData,
	})

}
