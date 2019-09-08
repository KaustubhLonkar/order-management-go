package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	product "github.appl.ge.com/geappliancesales/product-service-go/product"
	_ "github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	debug  = os.Getenv("DEBUG")
	rpcEnv = os.Getenv("RPCENV")
)

func main() {

	router := initRouter()
	router.Run(":8888")
}

func initRouter() *gin.Engine {

	if rpcEnv == "prod" || rpcEnv == "prd" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery(), plainLoggerWithWriter(gin.DefaultWriter))
	r.GET("/status", statusCheck)

	v1 := r.Group("/v1", v1Handler)

	v1.POST("addProduct", func(c *gin.Context) {
		product.AddProduct(c)
	})

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

func statusCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func exception(c *gin.Context) {
	c.JSON(500, gin.H{"success": false, "error": "Unable to process order"})
}

func v1Handler(c *gin.Context) {
	c.Set("version", 1)
	c.Next()
}
