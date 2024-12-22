package utils

import (
	"sync"
)

// Publisher 结构体表示发布者
type Publisher struct {
	subscribers sync.Map // 使用sync.Map作为订阅者存储
}

// NewPublisher 创建一个新的发布者实例
func NewPublisher() *Publisher {
	return &Publisher{
		subscribers: sync.Map{},
	}
}

// Subscribe 订阅消息
func (p *Publisher) Subscribe(sub string, fun func(any)) {
	value, ok := p.subscribers.Load(sub)
	if ok {
		value = append(value.([]func(any)), fun)
	} else {
		value = []func(any){fun}
	}
	p.subscribers.Store(sub, value)
}

// Unsubscribe 取消订阅
func (p *Publisher) Unsubscribe(sub string) {
	p.subscribers.Delete(sub)
}

// Publish 发布消息给订阅者
func (p *Publisher) Publish(topic string, data any) {
	value, ok := p.subscribers.Load(topic)
	if !ok {
		return
	}
	subs := value.([]func(interface{}))
	for _, fun := range subs {
		go fun(data)
	}
}

func main() {
	pub := NewPublisher()
	pub.Subscribe("test", func(data interface{}) {
		println(data.(string))
	})
	pub.Publish("test", "hello world")

}
