package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type handler struct {
	redisAddr string
	redisKey  string
	latency   time.Duration
	expiry    time.Duration
}

func (h *handler) HandleRequest(w http.ResponseWriter, req *http.Request) {
	// for sake of example, tear down redis connection after every request
	redisClient := redis.NewClient(&redis.Options{Addr: h.redisAddr})
	defer func() {
		cmd := redis.NewStringCmd("QUIT")
		redisClient.Process(cmd)
		redisClient.Close()
	}()

	if text, err := redisClient.Get(h.redisKey).Result(); err == nil {
		w.Write([]byte(text))
		return
	}

	time.Sleep(h.latency)

	text := "hello\n"
	expiry := time.Duration(rand.Int63n(int64(h.expiry)))
	redisClient.Set(h.redisKey, text, expiry)
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

	httpHandler := handler{
		redisAddr: *redisAddr,
		redisKey:  redisKey,
		latency:   *latency,
		expiry:    *expiry,
	}

	http.HandleFunc("/", httpHandler.HandleRequest)
	http.ListenAndServe(*addr, nil)
}
