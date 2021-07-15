package queue

import (
	"context"
	"email-dispatcher/config"
	"email-dispatcher/domain"
	"email-dispatcher/infra"
	"email-dispatcher/usecases"
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/segmentio/kafka-go"
)

// ReadMessageForDispatch reads the topic socket-service and calls
// the right action for the command
func ReadMessageForDispatch(group string, readTopic string) error {
	log.Printf("starting kafka consumer from topic %s", readTopic)
	brokers := strings.Split(config.GetConfig().KafkaBrokers(), ",")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Topic:   readTopic,
		GroupID: group,
		Brokers: brokers,
	})
	writer := &kafka.Writer{Addr: kafka.TCP(brokers...)}
	defer func() {
		log.Println("ending kafka function read with errors ")
		closeErr := reader.Close()
		if closeErr != nil {
			log.Println(closeErr)
		}
		closeErr = writer.Close()
		if closeErr != nil {
			log.Println(closeErr)
		}
	}()

	for {
		msg, msgErr := reader.ReadMessage(context.Background())
		if msgErr != nil {
			return msgErr
		}

		log.Printf("send email using amazon SES service ")

		awsConfig := aws.NewConfig().
			WithRegion(config.GetConfig().AmazonSESRegion())

		svc := infra.NewAmazonSES(awsConfig)
		task := usecases.NewEmailSender(svc)

		m := domain.Message{}

		JsonErr := json.Unmarshal(msg.Value, &m)
		if JsonErr != nil {
			log.Fatal("CANNOT UNMARSHAL: ", JsonErr)
		}

		task.SendEmail(m)
	}
}
