package core

import (
	"Day15/config"
	"Day15/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Getcart(c *gin.Context) {
	cartData := []models.Cart{}
	config.Db.First(&cartData, " user_id = ? ", c.GetUint("user_id"))

	c.JSON(200, gin.H{
		"status": "Success",
		"data":   cartData,
	})
}
func Addtocart(c *gin.Context) {

	if c.PostForm("qty") == "" || c.PostForm("sku") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Missing required parameter",
		})
		return
	}

	qt, _ := strconv.ParseInt(c.PostForm("qty"), 10, 64)
	cartData := models.Cart{
		Sku:     c.PostForm("sku"),
		Qty:     qt,
		User_id: c.GetUint("user_id"),
	}
	errz := config.Db.Create(&cartData)
	if errz.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"Message": errz.Error.Error(),
			"data":    cartData,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Ok",
			"Message": "Data Insert Cart",
			"data":    cartData,
		})

	}
}

func DeleteBySKU(c *gin.Context) {
	cartData := []models.Cart{}
	sku := c.Param("sku")
	config.Db.First(&cartData, "user_id=? and sku = ?", c.GetUint("user_id"), sku)
	//	Db.Find(&blogData)
	c.JSON(200, gin.H{
		"status": "Success",
		"data":   cartData,
	})

	errc := config.Db.Delete(&cartData)
	if errc.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"Message": errc.Error.Error(),
			"data":    cartData,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "Ok",
		"Message": "Data deleted Cart",
		"data":    cartData,
	})
}

func CheckoutCart(c *gin.Context) {
	cartData := []models.Cart{}
	transactionHeader := models.Transaksi{}
	transaction := []models.Transaksi_detail{}
	productdisplayArr := []models.Product_display{}

	config.Db.Find(&cartData, "user_id=? ", c.GetUint("user_id"))

	c.JSON(200, gin.H{
		"status": "Success",
		"data":   cartData,
	})

	var (
		total    float64
		subtotal float64
	)
	total = 0
	subtotal = 0

	transactionHeader.User_id = c.GetUint("user_id")
	config.Db.Save(&transactionHeader)
	idtransaksi := transactionHeader.ID

	for _, cartitem := range cartData {
		sku := cartitem.Sku
		qtycart := cartitem.Qty

		displayData := models.Product_display{}
		config.Db.First(&displayData, "sku = ?", sku)

		transactionItem := models.Transaksi_detail{}
		transactionItem.Id_transaksi = idtransaksi
		if displayData.Qty < cartitem.Qty {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Server Error",
				"Message": "refuse not available sku",
				"data":    displayData,
			})
			return
		}
		displayData.Qty = displayData.Qty - qtycart
		productdisplayArr = append(productdisplayArr, displayData)
		subtotal = (displayData.Price * float64(qtycart))
		total = total + subtotal
		transactionItem.Price = displayData.Price
		transactionItem.Qty = qtycart
		transactionItem.User_id = transactionHeader.User_id
		transactionItem.Sku = sku
		transactionItem.Subtotal = subtotal
		transaction = append(transaction, transactionItem)

	}
	config.Db.First(&transactionHeader, "ID=?", idtransaksi)
	transactionHeader.Total = total
	config.Db.Save(&transactionHeader)

	for _, dispitem := range productdisplayArr {
		config.Db.First(&dispitem, "sku = ?", dispitem.Sku)
		config.Db.Save(&dispitem)
		displayData := models.Product_display{}
		config.Db.First(&displayData, "sku = ?", dispitem.Sku)
		displayData.Qty = displayData.Qty - dispitem.Qty
		config.Db.Save(&displayData)
	}
	for _, transItem := range transaction {

		config.Db.Save(&transItem)
	}
	for _, cartitem := range cartData {
		errc := config.Db.Delete(&cartitem)
		if errc.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Server Error",
				"Message": errc.Error.Error(),
				"data":    cartData,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Ok",
		"Message": "Checkout success",
		"data":    cartData,
	})
}
func GetTransactionHistory(c *gin.Context) {
	transData := []models.Transaksi_detail{}
	config.Db.First(&transData, " user_id = ? ", c.GetUint("user_id"))

	c.JSON(200, gin.H{
		"status": "Success",
		"data":   transData,
	})
}
func GetAllReport(c *gin.Context) {
	filter, errz := strconv.ParseBool(c.Query("filter"))
	if errz != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"Message": "filter must be boolean",
			"data":    "filter is required",
		})
		return
	}
	dayStart := c.Query("dayStart")
	monthStart := c.Query("monthStart")
	yearStart := c.Query("yearStart")
	dayEnd := c.Query("dayEnd")
	monthEnd := c.Query("monthEnd")
	yearEnd := c.Query("yearEnd")

	transData := []models.Transaksi_detail{}
	if !filter {
		config.Db.Find(&transData)
	} else {
		if dayStart == "" || dayEnd == "" || monthStart == "" || monthEnd == "" || yearStart == "" || yearEnd == "" {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Server Error",
				"Message": "parameter day/month/year is required",
				"data":    "filter is required",
			})
			return
		}
		startD := yearStart + "-" + monthStart + "-" + dayStart
		endD := yearEnd + "-" + monthEnd + "-" + dayEnd
		config.Db.Find(&transData, "created_at between ? and ?", startD, endD)
	}

	c.JSON(200, gin.H{
		"status": "Success",
		"data":   transData,
	})
}
