package shipping

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/KaustubhLonkar/order-management-go/kafkaconfig"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

//GetShippingDetails func
func GetShippingDetails(c *gin.Context) {

	brokers := strings.Split(kafkaconfig.KafkaBrokerUrl, ",")
	config := kafka.ReaderConfig{
		Brokers:         brokers,
		GroupID:         kafkaconfig.KafkaClientId,
		Topic:           kafkaconfig.KafkaTopic,
		MinBytes:        10e3,            // 10KB
		MaxBytes:        10e6,            // 10MB
		MaxWait:         1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka.
		ReadLagInterval: -1,
	}

	reader := kafka.NewReader(config)
	defer reader.Close()

	//for {
	m, err := reader.ReadMessage(context.Background())
	if err != nil {
		log.Error().Msgf("error while receiving message: %s", err.Error())
		//		continue
		return
	}

	value := m.Value
	// if m.CompressionCodec == snappy.NewCompressionCodec() {
	// 	_, err = snappy.NewCompressionCodec().Decode(value, m.Value)
	// }

	// if err != nil {
	// 	log.Error().Msgf("error while receiving message: %s", err.Error())
	// 	continue
	// }

	fmt.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(value))
	kafkaconfig.Messages = append(kafkaconfig.Messages, string(value))
	//}

	c.JSON(200, kafkaconfig.Messages)

}
