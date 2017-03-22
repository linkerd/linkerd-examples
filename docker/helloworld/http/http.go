package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	text        string
	target      string
	podIp       string
	latency     time.Duration
	failureRate float64
	json        bool
}

func New(text, target, podIp string, latency time.Duration, failureRate float64, json bool) *Server {
	return &Server{
		text:        text,
		target:      target,
		podIp:       podIp,
		latency:     latency,
		failureRate: failureRate,
		json:        json,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/setLatency":
		s.handleLatency(w, req)
	case "/setFailureRate":
		s.handleFailureRate(w, req)
	default:
		s.handleRequest(w, req)
	}
}

func (s *Server) handleRequest(w http.ResponseWriter, req *http.Request) {
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

func (s *Server) handleLatency(w http.ResponseWriter, req *http.Request) {
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

func (s *Server) handleFailureRate(w http.ResponseWriter, req *http.Request) {
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

func (s *Server) callTarget(ctx *linkerdContext) (string, error) {
	req, err := http.NewRequest("GET", "http://"+s.target, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(ctx.withContext(req))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("invalid response %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

func (s *Server) writeJson(w http.ResponseWriter, text string) {
	jsonStr, err := json.Marshal(map[string]string{"api_result": text})
	if err != nil {
		http.Error(w, "error converting: "+text, http.StatusInternalServerError)
		return
	}

	w.Write(jsonStr)
}
