package kafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.ustock.cc/core/go-spring/types"

	"go.uber.org/zap"
	"sync"
	"time"
)

func StartConsumers(ctx context.Context , config KafkaConfig, handler Handler) {

	wg , ok := ctx.Value(sync.WaitGroup{}).(*sync.WaitGroup)
	channels := make([]chan *types.OpEvent, 0, config.ConsumerNum)
	var i int
	for i = 0; i < config.ConsumerNum; i++ {

		if ok {
			wg.Add(1)
		}
		ch := make(chan *types.OpEvent, 10)

		channels = append(channels, ch)

		go newReadLoop(ctx, ch, handler)

	}
}

func newReadLoop(ctx context.Context,  ch chan *types.OpEvent , handler Handler ) {

	wg,_ := ctx.Value(sync.WaitGroup{}).(*sync.WaitGroup)
	config , _ :=ctx.Value("kafka.config").(KafkaConfig)
	logger,_:= ctx.Value(zap.Logger{}).(*zap.Logger)
	defer wg.Done()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               config.KafkaAddrs,
		"group.id":                        config.ConsumerGroup,
		"session.timeout.ms":              30000,
		"enable.auto.commit":              false,
		"heartbeat.interval.ms":           3000,
		"go.application.rebalance.enable": true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": kafka.OffsetEnd.String()},
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to create consumer: %s", err))
	}
	defer c.Close()

	logger.Info("create a consumer ", zap.String("kafkaAddr",config.KafkaAddrs),
		zap.String("topic",config.Topic),
		zap.String("consumer",c.String()))

	var zapFields  = []zap.Field {zap.String("topic",config.Topic),
		zap.String("consumer",c.String())}

	topics := []string{config.Topic}
	if err := c.SubscribeTopics(topics, nil); err != nil {
		logger.Fatal("subscribe error" , append(zapFields, zap.String("error",err.Error()))...)
	}

	if _, _, err := c.QueryWatermarkOffsets(config.Topic, 0, 1000); err != nil {
		logger.Fatal("can't  QueryWatermarkOffsets" , append(zapFields, zap.String("error",err.Error()))...)
	}

	run := true
	err = handler.Start(ctx)
	if err != nil {
		logger.Fatal("Failed to start the  handler" , append(zapFields, zap.String("error",err.Error()))...)
	}

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	var tracedTickerIds []string
	msgs := make([]*kafka.Message, 0, config.ConsumerBatch)
	for run == true {
		select {
		case <-ctx.Done():
			logger.Info("the consumer reader loop :  receiving a exit signal ", zapFields...)
			run = false
		case <-ticker.C:
			if len(msgs) != 0 {
				err := handler.Process(ctx, msgs)

				if err != nil {
					logger.Info("failed to process ticker message " ,
						append(zapFields, zap.String("error",err.Error()))...)
				} else {
					last := msgs[len(msgs)-1]
					if _, err := c.CommitMessage(last); err != nil {


						logger.Info("failed to commit the ack msg" , append(zapFields, []zap.Field {
							zap.Int64("offset",int64(last.TopicPartition.Offset)),
							zap.String("error",err.Error()),
						}...)...)

					} else {
						//logger.Info("commit the message  to consumer" , append(zapFields, []zap.Field {
						//	zap.Int64("offset",int64(last.TopicPartition.Offset)),
						//}...)...)

						//logEntry.WithFields(logrus.Fields{"partition": last.TopicPartition.Partition, "offset": last.TopicPartition.Offset}).
						//	Debug("commit the message  to consumer ")
					}
					msgs=msgs[0:0]
				}
			}
		case traceEvent := <-ch:
			//logEntry.Infof("receiving a traceEvent %#v", traceEvent)
			switch traceEvent.Op {
			case types.OpTrace:
				tickerId := traceEvent.Body.(string)
				tracedTickerIds = append(tracedTickerIds, tickerId)
			case types.OpUntrace:
				tracedTickerIds = []string{}
			default:

				logger.Info("unsupported trace event" , append(zapFields, []zap.Field {
					zap.String("op",string(traceEvent.Op)),
				}...)...)
			}

		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:

				for _, tickerId := range tracedTickerIds {
					if string(e.Key) == tickerId {
						e.Opaque = true
						//logEntry.WithFields(logrus.Fields{"key": string(e.Key), "value": string(e.Value)}).
						//	WithField("traceTickerIds", tracedTickerIds).
						//	Infof("receiving a consumer message on partition:%d ", e.TopicPartition.Partition)
						break
					}
				}

				if e.Headers != nil {
					//logEntry.Infof("receiving a consumer message with headers: %v", e.Headers)
					handler.ProcessHeaders(ctx,e.Headers)
				}
				msgs = append(msgs, e)
				if len(msgs) >= config.ConsumerBatch {
					handler.Process(ctx, msgs) // 这里只是放入队列中，要么堵塞要么成功，不会失败,如果堵塞超过30秒会导致kafka rebalance
					if _, err := c.CommitMessage(e); err != nil {
						//logEntry.WithFields(logrus.Fields{"msg": e, "error": err}).Info("commit the msg fail")
					} else {
						//logEntry.WithFields(logrus.Fields{"partition": e.TopicPartition.Partition, "offset": e.TopicPartition.Offset}).
						//	Debug("commit the message  to consumer ")
					}
					msgs = msgs[0:0]
				}

			case kafka.AssignedPartitions:
				logger.Info("receiving a AssignedPartitions",zapFields...)
				c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				logger.Info("receiving a RevokedPartitions msg",zapFields...)
				c.Unassign()
			case kafka.PartitionEOF:
				logger.Info("Reached Partition EOF",zapFields...)
			case kafka.Error:
				logger.Info("Poll error" , append(zapFields, zap.String("error",e.String()))...)
			default:
				//logrus.Infof("Ignored %v", e)

				logger.Info("ignore the msg" , append(zapFields, zap.String("event",e.String()))...)

			}
		}
	}

	//接受方不应该关闭channel
	//close(ch)
	logger.Info("the consumer read loop exit ")

}


