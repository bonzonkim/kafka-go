package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/bonzonkim/kafka-go/config"
	"github.com/bonzonkim/kafka-go/kafka"
	"github.com/bonzonkim/kafka-go/logger"
	"go.uber.org/zap"
)

type ConsumerGroupHandler struct {
	clientID	 string
	Logger		 *zap.Logger
	consmuerGroup *kafka.MConsumerGroup
}
func (handler ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil 
}
func (handler ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil 
}

func (handler ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		ctx := handler.consmuerGroup.GetContextFromMsg(msg)
		fmt.Println("[Message Recieved] ", " timeStamp:", msg.Timestamp.Format("2006-01-02 15:04:05"), "consumerId:", handler.clientID, "context:", ctx, " - topic:", msg.Topic, " - key:", string(msg.Key), " - msgValue:", string(msg.Value), " - partition:", msg.Partition, " - offset:", msg.Offset)

		session.MarkMessage(msg, "")
	}
	return nil
}

func startConsumer(ctx context.Context, cfg *config.Config, log *zap.Logger) error {
	clientID := fmt.Sprintf("consumer-%d", rand.Intn(1000))

	group, err := kafka.NewConsumerGroup(&cfg.Kafka, cfg.Kafka.Topic, "my-consumer-group", clientID, log)
	if err != nil {
		return err
	}
	defer group.Close()

	handler := ConsumerGroupHandler{
		Logger:	log,
		clientID: clientID,
		consmuerGroup: group,
	}
	group.RegisterHandlerAndConsumeMessages(ctx, handler)

	return nil
}

func main() {
	cfg, _ := config.LoadConfig("config/config.yaml")
	log := logger.NewLogger("consumer", cfg.Log.RotationSize, cfg.Log.RotationCount)

	// Start Consumer in Background
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := startConsumer(ctx, cfg, log); err != nil {
			log.Fatal("Failed to start consumer", zap.Error(err))
		}
	}()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	cancel()
	time.Sleep(2 * time.Second)
	log.Info("Shutting down gracefully")
}
