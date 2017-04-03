package redis

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

const expiration time.Duration = time.Minute

type Client struct {
	client *redis.Client
	key    string
}

func New(addr string) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	key := fmt.Sprintf("proc:%d", os.Getpid())
	return &Client{client: client, key: key}
}

func (c *Client) Set(val string) error {
	return c.client.Set(c.key, val, expiration).Err()
}

func (c *Client) Get() (string, error) {
	return c.client.Get(c.key).Result()
}
