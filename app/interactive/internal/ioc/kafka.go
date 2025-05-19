package ioc

import (
	"github.com/IBM/sarama"
	saramax "github.com/LEILEI0628/GinPro/Saramax"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/conf"
	"github.com/LEILEI0628/GoWeb-MicroServices/app/interactive/internal/events"
)

func InitKafka(conf *conf.Bootstrap) sarama.Client {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	client, err := sarama.NewClient(conf.GetMessageQueue().GetKafka().GetAddrs(), saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

// NewConsumers 面临的问题依旧是所有的 Consumer 在这里注册一下
func NewConsumers(c1 *events.InteractiveReadEventConsumer) []saramax.Consumer {
	return []saramax.Consumer{c1}
}
