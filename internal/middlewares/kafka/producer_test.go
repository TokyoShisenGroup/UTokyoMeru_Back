package kafka

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewKafkaProducer(t *testing.T) {
	t.Log("开始测试 Kafka Producer 创建...")

	broker := "43.133.12.107:30092"
	topic := "myTopic"

	t.Logf("尝试连接到 Kafka broker: %s, topic: %s", broker, topic)

	// 初始化 Kafka producer
	producer, err := NewProducer(broker, topic)
	if err != nil {
		t.Errorf("创建 Producer 失败: %v", err)
		t.FailNow()
	}

	assert.NoError(t, err, "Producer 创建应该没有错误")
	assert.NotNil(t, producer, "Producer 不应该为 nil")

	t.Log("Producer 创建成功，尝试发送测试消息...")

	// 尝试发送测试消息
	testMessage := "这是一条测试消息"
	err = producer.ProduceMessage([]byte(testMessage))
	if err != nil {
		t.Errorf("发送消息失败: %v", err)
	} else {
		t.Log("测试消息发送成功")
	}

	// 等待一小段时间确保消息被处理
	time.Sleep(2 * time.Second)

	t.Log("测试完成，准备清理资源...")

	// 清理资源
	producer.P.Close()
	t.Log("Producer 已关闭，测试结束")

	// 打印 Producer 的详细信息
	t.Logf("Producer 详细信息: %+v", producer)
}
