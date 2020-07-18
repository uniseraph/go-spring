package kafka


type KafkaConfig struct {
	KafkaAddrs    string
	Topic         string
	ConsumerNum   int
	ConsumerBatch int
	ConsumerGroup string
}
