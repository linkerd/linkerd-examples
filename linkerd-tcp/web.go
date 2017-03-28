package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

type redisClient struct {
	addr   string
	key    string
	expiry time.Duration
	client *redis.Client
	sync.RWMutex
}

type handler struct {
	redisClient *redisClient
	latency     time.Duration
}

func (c *redisClient) Get() (string, error) {
	c.RLock()
	defer c.RUnlock()
	return c.client.Get(c.key).Result()
}

func (c *redisClient) Set(text string) error {
	c.RLock()
	defer c.RUnlock()
	expiry := time.Duration(rand.Int63n(int64(c.expiry)))
	return c.client.Set(c.key, text, expiry).Err()
}

func (c *redisClient) Refresh() {
	c.Lock()
	defer c.Unlock()
	if c.client != nil {
		c.client.Close()
	}
	c.client = redis.NewClient(&redis.Options{Addr: c.addr})
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if text, err := h.redisClient.Get(); err == nil {
		w.Write([]byte(text))
		return
	}

	time.Sleep(h.latency)

	text := "hello\n"
	h.redisClient.Set(text)
	w.Write([]byte(text))
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	addr := flag.String("addr", ":7770", "service address to serve on")
	redisAddr := flag.String("redis-addr", "127.0.0.1:6379", "address of redis cluster")
	latency := flag.Duration("latency", 300*time.Millisecond, "time to sleep on cache miss")
	expiry := flag.Duration("expiry", 3*time.Minute, "max cache key expire time")
	flag.Parse()

	redisKey := "app:" + strconv.FormatInt(rand.Int63(), 16)
	fmt.Printf("serving on %s, caching on %s\n", *addr, redisKey)

	client := redisClient{
		addr:   *redisAddr,
		key:    redisKey,
		expiry: *expiry,
	}
	client.Refresh()

	// refresh connection every 5 seconds, with jitter
	jitter := time.Duration(rand.Int63n(int64(*latency)))
	go func() {
		for _ = range time.Tick(5*time.Second + jitter) {
			client.Refresh()
		}
	}()

	httpServer := &http.Server{
		Addr:        *addr,
		ReadTimeout: 10 * time.Second,
		Handler: &handler{
			redisClient: &client,
			latency:     *latency,
		},
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
