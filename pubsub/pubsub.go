package pubsub

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewPubSub() *PubSub {
	return &PubSub{context.Background(), redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})}
}

var pubSub *PubSub

func GetPubSub() *PubSub {
	if pubSub == nil {
		pubSub = NewPubSub()
	}
	return pubSub
}

type PubSub struct {
	ctx context.Context
	rdb *redis.Client
}

func (p *PubSub) Ctx() context.Context {
	return p.ctx
}

func (p *PubSub) Publish(channel string, message interface{}) {
	fmt.Println(channel, message)
	p.rdb.Publish(p.ctx, channel, message)
}

func (p *PubSub) Subscribe(channel string) *redis.PubSub {
	fmt.Println(channel)
	return p.rdb.Subscribe(p.ctx, channel)
}
