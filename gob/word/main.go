package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
)

type (
	WordSvc struct {
		words []string
	}

	wordRsp struct{ Word string }
)

// the wordsvc api picks a cool word at random
func (svc *WordSvc) ServeHTTP(rspw http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		switch req.Method {
		case "GET":
			word := svc.randomWord()
			if word == "" {
				rspw.WriteHeader(http.StatusInternalServerError)
				return
			}
			body, err := json.Marshal(&wordRsp{word})
			if err != nil {
				rspw.WriteHeader(http.StatusInternalServerError)
				return
			}

			rspw.WriteHeader(http.StatusOK)
			rspw.Write(body)
			return

		default:
			rspw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

	default:
		rspw.WriteHeader(http.StatusNotFound)
		return
	}
}

func (svc *WordSvc) randomWord() string {
	n := len(svc.words)
	switch n {
	case 0:
		return ""
	default:
		idx := rand.Int() % n
		return svc.words[idx]
	}
}

func dieIf(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %s. Try --help for help.\n", err)
	os.Exit(-1)
}

func main() {
	//
	// Parse flags
	//
	addr := flag.String("srv", ":8282", "TCP address to listen on (in host:port form)")
	flag.Parse()
	if flag.NArg() != 0 {
		dieIf(fmt.Errorf("expecting zero arguments but got %d", flag.NArg()))
	}

	//
	// Setup Http server
	//
	server := &http.Server{
		Addr: *addr,
		Handler: &WordSvc{
			words: []string{
				"banana",
				"bees",
				"cmon",
				"gob",
				"hermano",
				"illusion",
				"same",
				"",
				"",
			},
		},
	}
	dieIf(server.ListenAndServe())
}
