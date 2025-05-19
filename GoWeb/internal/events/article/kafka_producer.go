package events

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(producer sarama.SyncProducer) Producer {
	return &KafkaProducer{producer: producer}
}

func (k *KafkaProducer) ProduceReadEvent(ctx context.Context, evt ReadEvent) error {
	data, err := json.Marshal(evt)
	if err != nil {
		return err
	}
	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{ // 没有（无法）传context
		Topic: "read_article", // 可以使用外部传输
		Value: sarama.ByteEncoder(data),
	})
	// 如果重试逻辑很简单，放这里；复杂的话用装饰器
	return err
}
