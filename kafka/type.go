package kafka

import "context"

type KafkaConfig struct {
	KafkaAddrs    string
	Topic         string
	ConsumerNum   int
	ConsumerBatch int
	ConsumerGroup string
}

type kafkaConfigContextKey struct {}

func NewContext(ctx context.Context, cfg *KafkaConfig) context.Context {
	return context.WithValue(ctx, kafkaConfigContextKey{},cfg)
}

func FromContext(ctx context.Context) (*KafkaConfig,bool){
	v,ok := ctx.Value(kafkaConfigContextKey{}).(*KafkaConfig)
	return v,ok
}