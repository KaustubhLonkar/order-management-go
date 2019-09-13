package products

import (
	"net/http"
	"strconv"

	model "github.com/KaustubhLonkar/order-management-go/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type product struct {
	gorm.Model
	ProductName     string `json:"product_name"`
	ProductQuantity int    `json:"product_quantity"`
	ProductPrice    int    `json:"product_price"`
}

//TransformedProduct struct
type TransformedProduct struct {
	ProductID       uint   `json:"id"`
	ProductName     string `json:"product_name"`
	ProductQuantity int    `json:"product_quantity"`
	ProductPrice    int    `json:"product_price"`
}

func addProductDB(prod *product) {
	model.Db.AutoMigrate(&prod)
}

// AddProduct func
func AddProduct(c *gin.Context) {

	price, _ := strconv.Atoi(c.PostForm("pprice"))
	quantity, _ := strconv.Atoi(c.PostForm("pquantity"))
	prod := product{

		ProductName:     c.PostForm("pname"),
		ProductPrice:    price,
		ProductQuantity: quantity,
	}
	addProductDB(&prod)
	model.Db.Save(&prod)
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Product added successfully",
	})

}

// GetProducts func
func GetProducts(c *gin.Context) {
	var products []product
	var _products []TransformedProduct

	model.Db.Find(&products)
	if len(products) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No products found !!"})
		return
	}
	for _, item := range products {
		_products = append(_products, TransformedProduct{ProductID: item.ID, ProductName: item.ProductName, ProductQuantity: item.ProductQuantity, ProductPrice: item.ProductPrice})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   _products,
	})
}

//IsProductAvailable to check availability
func IsProductAvailable(pname string, qty int) (bool, product, int) {
	available := false
	amt := 0
	var prod product
	model.Db.Where("product_name = ?", pname).First(&prod)

	if pname == prod.ProductName && qty <= prod.ProductQuantity {
		available = true
		amt = prod.ProductPrice * qty
		remainingQuantity := prod.ProductQuantity - qty
		model.Db.Model(&prod).Update("product_quantity", remainingQuantity)
		return available, prod, amt
	}
	return available, prod, amt
}
