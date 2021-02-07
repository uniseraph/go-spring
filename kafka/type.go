package kafka

import (
	"context"
	"sync"
)

type KafkaConfig struct {
	KafkaAddrs    string
	Topic         string
	ConsumerNum   int
	ConsumerBatch int
	ConsumerGroup string
}

type waitGroupKey struct {}

func NewContext(ctx context.Context, group *sync.WaitGroup) context.Context {
	return context.WithValue(ctx, waitGroupKey{},group)
}

func FromContext(ctx context.Context) (*sync.WaitGroup,bool){
	v,ok := ctx.Value(waitGroupKey{}).(*sync.WaitGroup)
	return v,ok
}