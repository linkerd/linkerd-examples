package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type server struct {
	text        string
	target      string
	podIp       string
	latency     time.Duration
	failureRate float64
	json        bool
}

func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/setLatency":
		s.handleLatency(w, req)
	case "/setFailureRate":
		s.handleFailureRate(w, req)
	default:
		s.handleRequest(w, req)
	}
}

func (s *server) handleRequest(w http.ResponseWriter, req *http.Request) {
	time.Sleep(s.latency)
	if rand.Float64() < s.failureRate {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	text := s.text
	if s.podIp != "" {
		text += fmt.Sprintf(" (%s)", s.podIp)
	}

	if s.target != "" {
		targetText, err := s.callTarget(getContext(req))
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		text += fmt.Sprintf(" %s", targetText)
	}

	if s.json {
		s.writeJson(w, text)
	} else {
		w.Write([]byte(text + "!"))
	}
}

func (s *server) handleLatency(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	latencyParam := req.URL.Query().Get("latency")
	if latencyParam == "" {
		http.Error(w, "missing required parameter: latency", http.StatusBadRequest)
		return
	}

	latency, err := time.ParseDuration(latencyParam)
	if err != nil {
		http.Error(w, "latency is not a valid duration", http.StatusBadRequest)
		return
	}

	s.latency = latency
	fmt.Println("set latency to", latency)
	w.Write([]byte("ok"))
}

func (s *server) handleFailureRate(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	failureRateParam := req.URL.Query().Get("failureRate")
	if failureRateParam == "" {
		http.Error(w, "missing required parameter: failureRate", http.StatusBadRequest)
		return
	}

	failureRate, err := strconv.ParseFloat(failureRateParam, 64)
	if err != nil {
		http.Error(w, "failureRate is not a valid float", http.StatusBadRequest)
		return
	}

	if failureRate < 0.0 || failureRate > 1.0 {
		http.Error(w, "failureRate must be between 0.0 and 1.0", http.StatusBadRequest)
		return
	}

	s.failureRate = failureRate
	fmt.Println("set failure rate to", failureRate)
	w.Write([]byte("ok"))
}

func (s *server) callTarget(ctx *linkerdContext) ([]byte, error) {
	req, err := http.NewRequest("GET", "http://"+s.target, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(ctx.withContext(req))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (s *server) writeJson(w http.ResponseWriter, text string) {
	jsonStr, err := json.Marshal(map[string]string{"api_result": text})
	if err != nil {
		http.Error(w, "error converting: "+text, http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}

type linkerdContext map[string]string

func getContext(req *http.Request) *linkerdContext {
	ctx := make(linkerdContext)
	for key, _ := range req.Header {
		if strings.HasPrefix(strings.ToLower(key), "l5d-ctx") {
			ctx[key] = req.Header.Get(key)
		}
	}
	return &ctx
}

func (lc *linkerdContext) withContext(req *http.Request) *http.Request {
	req2 := new(http.Request)
	*req2 = *req
	for key, val := range *lc {
		req2.Header.Set(key, val)
	}
	return req2
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	addr := flag.String("addr", ":7777", "address to serve on")
	text := flag.String("text", "Hello", "text to serve")
	target := flag.String("target", "", "target service to call before returning")
	latency := flag.Duration("latency", 0, "time to sleep before processing request")
	failureRate := flag.Float64("failure-rate", 0.0, "rate of 500 responses to return")
	json := flag.Bool("json", false, "return json instead of plaintext responses")
	flag.Parse()

	fmt.Println("starting http server on", *addr)

	serverText := *text
	if envText := os.Getenv("TARGET_WORLD"); envText != "" {
		serverText = envText
	}

	httpServer := &http.Server{
		Addr:         *addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler: &server{
			text:        serverText,
			target:      *target,
			podIp:       os.Getenv("POD_IP"),
			latency:     *latency,
			failureRate: *failureRate,
			json:        *json,
		},
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
