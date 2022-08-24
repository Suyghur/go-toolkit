//@File     producer.go
//@Time     2022/08/23
//@Author   #Suyghur,

package ckafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/suyghur/go-toolkit/trace"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"strings"
)

type ProducerConfig struct {
	Brokers   []string `json:""`
	Topic     string   `json:""`
	User      string   `json:",optional"`
	Passwd    string   `json:",optional"`
	Partition int      `json:",default=1"`
}

type Producer struct {
	topic         string
	brokers       []string
	config        ProducerConfig
	producer      *kafka.Producer
	msgChan       chan kafka.Event
	partitionIds  []int32
	nextPartition int
}

func MustNewProducer(config ProducerConfig) *Producer {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(config.Brokers, ","),
	}
	if config.User != "" {
		_ = configMap.Set("sasl.mechanisms=PLAIN")
		_ = configMap.Set("security.protocol=SASL_PLAINTEXT")
		_ = configMap.Set("sasl.username=" + config.User)
		_ = configMap.Set("sasl.password=" + config.Passwd)
		// 超时时间
		_ = configMap.Set("request.timeout.ms=500")
		// 重试
		_ = configMap.Set("retries=0")
	}
	producer, err := kafka.NewProducer(configMap)
	if err != nil {
		logx.Errorf("kafka.NewProducer error: %v", err)
		panic(err)
	}
	admin, err := kafka.NewAdminClientFromProducer(producer)
	if err != nil {
		logx.Errorf("kafka.NewAdminClientFromProducer error: %v", err)
		panic(err)
	}
	var partitionIds []int32
	metadata, err := admin.GetMetadata(&config.Topic, false, 5000)
	if err != nil {
		logx.Errorf("kafka.GetMetadata error: %v", err)
		for i := 0; i < config.Partition; i++ {
			partitionIds = append(partitionIds, int32(i))
		}
	} else {
		for i, partition := range metadata.Topics[config.Topic].Partitions {
			logx.Infof("kafka.GetMetadata: %s:%d partition: %+v", config.Topic, i, partition)
			partitionIds = append(partitionIds, partition.ID)
		}
	}
	ch := make(chan kafka.Event)
	return &Producer{
		topic:         config.Topic,
		brokers:       config.Brokers,
		config:        config,
		producer:      producer,
		msgChan:       ch,
		partitionIds:  partitionIds,
		nextPartition: 0,
	}
}

func (p *Producer) SendMessage(ctx context.Context, bMsg []byte, key string) (partition int32, offset int64, err error) {
	defer func() {
		p.nextPartition++
		if p.nextPartition >= len(p.partitionIds) {
			p.nextPartition = 0
		}
	}()
	for i := 0; i < 3; i++ {
		select {
		case <-ctx.Done():
			return 0, 0, ctx.Err()
		default:
			trace.StartFuncSpan(ctx, "kafka.Producer.SendMessage:key:"+key+"retrycount:"+strconv.Itoa(i), func(ctx context.Context) {
				err = p.producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic:     &p.topic,
						Partition: p.partitionIds[p.nextPartition],
					}, Key: []byte(key), Value: bMsg,
				}, p.msgChan)
			})
			if err != nil {
				continue
			}
			e := <-p.msgChan
			ev := e.(*kafka.Message)
			if ev.TopicPartition.Error != nil {
				err = ev.TopicPartition.Error
				continue
			}
			partition, offset = ev.TopicPartition.Partition, int64(ev.TopicPartition.Offset)
			return
		}
	}
	return
}
