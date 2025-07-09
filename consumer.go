package main

import (
	"sync/atomic"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/ha"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

type Consumer struct {
	c            *ha.ReliableConsumer
	consumed     atomic.Uint64
	totalEntries atomic.Uint64
}

func NewConsumer(env *stream.Environment, name string, opts *stream.ConsumerOptions) (c *Consumer, err error) {

	c = &Consumer{}

	rc, err := ha.NewReliableConsumer(env, name, opts, c.handleMessage)
	if err != nil {
		return nil, err
	}

	c.c = rc

	return c, nil
}

func (c *Consumer) handleMessage(consumerContext stream.ConsumerContext, message *amqp.Message) {
	c.consumed.Add(1)
	c.totalEntries.Add(uint64(consumerContext.GetEntriesCount()))
}
