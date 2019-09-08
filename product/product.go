package product

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//product struct
type inventory struct {
	gorm.Model
	ProductName     string `json:"product_name"`
	ProductQuantity int    `json:"product_quantity"`
	ProductPrice    int    `json:"product_price"`
}

var db *gorm.DB

func init() {
	var err error

	db, err = gorm.Open("mysql", "root:root@/OrderManagement?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&inventory{})
}

//AddProduct func
func AddProduct(c *gin.Context) {

	price, _ := strconv.Atoi(c.PostForm("pprice"))
	quantity, _ := strconv.Atoi(c.PostForm("pquantity"))
	prod := inventory{
		ProductName:     c.PostForm("pname"),
		ProductPrice:    price,
		ProductQuantity: quantity,
	}
	db.Save(&prod)
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Product added successfully",
	})
}
