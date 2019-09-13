package order

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/KaustubhLonkar/order-management-go/kafkaconfig"
	model "github.com/KaustubhLonkar/order-management-go/model"
	products "github.com/KaustubhLonkar/order-management-go/products"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

type order struct {
	gorm.Model
	Product     products.TransformedProduct `gorm:"foreignkey:json:"id"`
	OrderAmount int                         `json:"order_amount"`
}

var (
	logger = log.With().Str("pkg", "main").Logger()
)

type customer struct {
	Name    string
	Address string
	PinCode string
}

var custInfo = customer{
	Name:    "shon",
	Address: "Bellandur",
	PinCode: "560103",
}

func addOrderDB(ord *order) {
	model.Db.AutoMigrate(&ord)
}

// PlaceOrder func
func PlaceOrder(c *gin.Context) {
	pname := c.PostForm("pname")
	qty, _ := strconv.Atoi(c.PostForm("pqty"))

	avail, prod, amt := products.IsProductAvailable(pname, qty)
	if avail {

		fmt.Println("placed order is succefull...")
		ord := order{

			Product: products.TransformedProduct{
				ProductID:       prod.ID,
				ProductName:     prod.ProductName,
				ProductQuantity: qty,
				ProductPrice:    prod.ProductPrice,
			},
			OrderAmount: amt,
		}
		ord.Product.ProductID = prod.ID
		//fmt.Println(ord)
		addOrderDB(&ord)
		model.Db.Save(&ord)
		kafkaProducer, err := kafkaconfig.Configure(strings.Split(kafkaconfig.KafkaBrokerUrl, ","), kafkaconfig.KafkaClientId, kafkaconfig.KafkaTopic)
		if err != nil {
			logger.Error().Str("error", err.Error()).Msg("unable to configure kafka")
			return
		}
		defer kafkaProducer.Close()

		msg := kafkaconfig.Message{
			OrderID:         ord.ID,
			ProductID:       ord.Product.ProductID,
			ProductName:     ord.Product.ProductName,
			ProductQuantity: ord.Product.ProductQuantity,
			TotalAmount:     ord.OrderAmount,
			CustomerName:    custInfo.Name,
			CustomerAddress: custInfo.Address,
			CustomerPinCode: custInfo.PinCode,
		}
		kafkaconfig.PostDataToKafka(c, msg)
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "placed order is successfull...!",
			"data":    ord,
		})
		return
	} else {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Product found !!"})
		return
	}

}
