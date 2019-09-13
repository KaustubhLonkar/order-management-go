package kafkaconfig

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	//kafka "github.com/yusufsyaifudin/go-kafka-example/dep/kafka"

	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

var (
	listenAddrApi  string
	KafkaBrokerUrl string
	kafkaVerbose   bool
	KafkaClientId  string
	KafkaTopic     string
	Messages       []string
)

func init() {
	//flag.StringVar(&listenAddrApi, "listen-address", "0.0.0.0:9000", "Listen address for api")
	flag.StringVar(&KafkaBrokerUrl, "kafka-brokers", "localhost:9092", "Kafka brokers in comma separated value")
	flag.BoolVar(&kafkaVerbose, "kafka-verbose", true, "Kafka verbose logging")
	flag.StringVar(&KafkaClientId, "kafka-client-id", "my-kafka-client", "Kafka client id to connect")
	flag.StringVar(&KafkaTopic, "kafka-topic", "sql-insert", "Kafka topic to push") //sql-insert is a topic name
	flag.Parse()

}

// Message
type Message struct {
	OrderID         uint
	ProductID       uint
	ProductName     string
	ProductQuantity int
	TotalAmount     int
	CustomerName    string
	CustomerAddress string
	CustomerPinCode string
}

var writer *kafka.Writer

// Configure func
func Configure(KafkaBrokerUrls []string, clientId string, topic string) (w *kafka.Writer, err error) {
	fmt.Println("Started kafka configuration")
	// flag.StringVar(&listenAddrApi, "listen-address", "0.0.0.0:9000", "Listen address for api")
	// flag.StringVar(&KafkaBrokerUrl, "kafka-brokers", "localhost:9092,localhost:9093,localhost:9094,localhost:9095", "Kafka brokers in comma separated value")
	// flag.BoolVar(&kafkaVerbose, "kafka-verbose", true, "Kafka verbose logging")
	// flag.StringVar(&KafkaClientId, "kafka-client-id", "my-kafka-client", "Kafka client id to connect")
	// flag.StringVar(&KafkaTopic, "kafka-topic", "foo", "Kafka topic to push")
	// flag.Parse()

	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientId,
	}
	config := kafka.WriterConfig{
		Brokers:          KafkaBrokerUrls,
		Topic:            topic,
		Balancer:         &kafka.LeastBytes{},
		Dialer:           dialer,
		WriteTimeout:     10 * time.Second,
		ReadTimeout:      10 * time.Second,
		CompressionCodec: snappy.NewCompressionCodec(),
	}
	w = kafka.NewWriter(config)
	writer = w
	return w, nil
}

//Push
func Push(parent context.Context, key, value []byte) (err error) {
	message := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	return writer.WriteMessages(parent, message)
}

func PostDataToKafka(ctx *gin.Context, message Message) {
	parent := context.Background()
	defer parent.Done()

	ctx.Bind(message)
	formInBytes, err := json.Marshal(message)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("error while marshalling json: %s", err.Error()),
			},
		})

		ctx.Abort()
		return
	}

	err = Push(parent, nil, formInBytes)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": map[string]interface{}{
				"message": fmt.Sprintf("error while push message into kafka: %s", err.Error()),
			},
		})

		ctx.Abort()
		return
	}

}
