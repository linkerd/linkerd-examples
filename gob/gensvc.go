package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type (
	GenSvc struct{}

	genReq struct {
		Text  string
		Limit uint
	}
)

// gensvc web api
func (svc *GenSvc) ServeHTTP(rspw http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		switch req.Method {
		case "POST":
			var gen genReq
			if err := json.NewDecoder(req.Body).Decode(&gen); err != nil {
				rspw.WriteHeader(http.StatusBadRequest)
				return
			}
			if err := svc.generate(gen.Text, gen.Limit, rspw); err != nil {
				rspw.WriteHeader(http.StatusInternalServerError)
				return
			}
			return

		default:
			rspw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

	case "/health":
		switch req.Method {
		case "GET":
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

// Writes to the stream until `limit` writes have been completed or
// the stream is closed (i.e. because the client disconnects).
func (svc *GenSvc) generate(text string, limit uint, stream io.Writer) error {
	if _, err := stream.Write([]byte(text)); err != nil {
		return err
	}
	doWrite := func() bool {
		_, err := stream.Write([]byte(" " + text))
		return err == nil
	}
	if limit == 0 {
		for {
			if ok := doWrite(); !ok {
				break
			}
		}
	} else {
		for i := uint(0); i != limit; i++ {
			if ok := doWrite(); !ok {
				break
			}
		}
	}
	return nil
}

func dieIf(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %s. Try --help for help.\n", err)
	os.Exit(-1)
}

func main() {
	addr := flag.String("srv", ":8181", "TCP address to listen on (in host:port form)")
	flag.Parse()
	if flag.NArg() != 0 {
		dieIf(fmt.Errorf("expecting zero arguments but got %d", flag.NArg()))
	}
	server := &http.Server{
		Addr:    *addr,
		Handler: &GenSvc{},
	}
	dieIf(server.ListenAndServe())
}
