package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	order "github.com/KaustubhLonkar/order-management-go/order"
	products "github.com/KaustubhLonkar/order-management-go/products"
	shipping "github.com/KaustubhLonkar/order-management-go/shipping"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// used to match which service is being called
const (
	ADDPRODUCT  = "/addProduct"
	GETPRODUCTS = "/getProducts"
	PLACEORDER  = "/placeOrder"
	SHIPPING    = "/getShipping"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	router := initRouter()
	router.Run(":8888")
}

func initRouter() *gin.Engine {

	r := gin.New()
	p := ginprometheus.NewPrometheus("gin")
	p.Use(r)
	r.Use(static.Serve("/", static.LocalFile("./view", true)))
	r.Use(gin.Recovery(), plainLoggerWithWriter(gin.DefaultWriter))

	r.GET("/status", statusCheck)
	r.POST("/addProduct", requestRouter)
	r.GET("/getProducts", requestRouter)
	r.POST("/placeOrder", requestRouter)

	return r
}

// PlainLoggerWithWriter mimics the Gin LoggerWithWriter without the colors
func plainLoggerWithWriter(out io.Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		if c.Request.URL.Path != "/status" {
			fmt.Fprintf(out, "%s [%s] %s [%v] \"%s %s %s\" %d %d %v %s %s %s \"%s\"\n",
				c.ClientIP(),
				c.Request.UserAgent(),
				c.Request.Header.Get(gin.AuthUserKey),
				end.Format("02/Jan/2006:15:04:05 -0700"),
				c.Request.Method,
				c.Request.URL.Path,
				c.Request.Proto,
				c.Writer.Status(),
				c.Writer.Size(),
				fmt.Sprintf("%.4f", latency.Seconds()),
				c.Request.Header.Get("RequestType"),
				c.Request.Header.Get("ResponseSource"),
				c.Request.Form.Encode(),
				c.Request.Header.Get("ResponseBody"),
			)
		}
	}
}

// statusCheck returns a 200/OK when called if we can contact the be env
func statusCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func exception(c *gin.Context) {
	c.JSON(500, gin.H{"success": false, "error": "Unable to process order"})
}

func requestRouter(c *gin.Context) {

	path := c.Request.URL.Path
	fmt.Println("The obtained path is:- ", path)
	switch path {
	case ADDPRODUCT:
		products.AddProduct(c)
	case GETPRODUCTS:
		products.GetProducts(c)
	case PLACEORDER:
		order.PlaceOrder(c)
	case SHIPPING:
		shipping.GetShippingDetails(c)
	}
}
