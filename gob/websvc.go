package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type (
	GobWeb struct {
		genSvc  GenSvc
		wordSvc WordSvc
	}

	// A generic client for a downstream service
	Client struct {
		Name      string
		dstScheme string
		dstAddr   string
		client    *http.Client
	}

	contextHeaders map[string][]string

	wordSvc struct{ Client }
	WordSvc interface {
		Word(ctx contextHeaders) (string, error)
	}

	genSvc struct{ Client }
	GenSvc interface {
		Gen(ctx contextHeaders, text string, limit uint) (io.ReadCloser, error)
	}
	// Used
	helpCtx struct {
		Host string
	}

	genReq struct {
		Text  string
		Limit uint
	}

	wordRsp struct{ Word string }
)

// Default text for / and /help endpoints
var HelpTemplate = template.Must(template.New("help").Parse(`Gob's web service!

Send me a request like:

  {{.Host}}/gob

You can tell me what to say with:

  {{.Host}}/gob?text=WHAT_TO_SAY&limit=NUMBER
`))

// Gob's Web service
func (gob *GobWeb) ServeHTTP(rspw http.ResponseWriter, req *http.Request) {
	var err error
	ctx := getContextHeaders(req)

	switch req.URL.Path {
	case "/", "/help":
		switch req.Method {
		case "GET":
			rspw.Header().Set("content-type", "text/plain")
			if err = HelpTemplate.Execute(rspw, &helpCtx{req.Host}); err != nil {
				fmt.Fprintf(os.Stderr, "template error: %s", err.Error())
				return
			}
			return

		default:
			rspw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

	case "/gob":
		switch req.Method {
		case "GET":

			params := req.URL.Query()
			text := params.Get("text")
			if text == "" {
				text, err = gob.wordSvc.Word(ctx)
				if err != nil {
					fmt.Println("wordsvc error: " + err.Error())
					rspw.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			if text == "" {
				fmt.Println("could not load text")
				rspw.WriteHeader(http.StatusInternalServerError)
				return
			}

			limit := uint(0)
			limitstr := params.Get("limit")
			if limitstr != "" {
				limit64, err := strconv.ParseUint(limitstr, 10, 32)
				if err != nil {
					rspw.WriteHeader(http.StatusBadRequest)
					return
				}
				limit = uint(limit64)
			}

			stream, err := gob.genSvc.Gen(ctx, text, limit)
			if err != nil {
				fmt.Println("error generating: " + err.Error())
				rspw.WriteHeader(http.StatusInternalServerError)
				return
			}

			rspw.Header().Set("content-type", "text/plain")
			rspw.WriteHeader(http.StatusOK)
			io.Copy(rspw, stream)
			stream.Close()
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

//
// Client commonalities
//

// Client helper for sending downstream requests:
// - overrides the HTTP scheme
// - overrides the destination address
// - overrides the Host header
// - sets context headers (for routing and tracing)
func (svc *Client) request(
	ctx contextHeaders,
	method string,
	u *url.URL,
	body io.Reader,
) (*http.Response, error) {
	u.Scheme = svc.dstScheme
	u.Host = svc.dstAddr
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	ctx.applyTo(req)
	req.Host = svc.Name
	return svc.client.Do(req)
}

//
// XXX in the real word this would be a `context.Context`, but we're
// minimizing dependencies for this example.
//

// extract router-specific headers to support dynamic routing & tracing
func getContextHeaders(req *http.Request) contextHeaders {
	headers := make(map[string][]string)
	for k, values := range req.Header {
		prefixed := func(s string) bool { return strings.HasPrefix(k, s) }
		if prefixed("L5d-") || prefixed("Dtab-") || prefixed("X-Dtab-") {
			headers[k] = values
		}
	}
	return headers
}

func (ctx contextHeaders) applyTo(req *http.Request) {
	for k, values := range ctx {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}
}

//
// gensvc: generates a stream of data given a word

// Make a client to gensvc
func MakeGenSvc(scheme, addr string) GenSvc {
	return &genSvc{Client{"gen", scheme, addr, &http.Client{}}}
}

func (svc *genSvc) Gen(
	ctx contextHeaders,
	text string,
	limit uint,
) (io.ReadCloser, error) {
	if text == "" {
		return nil, errors.New("no text specified")
	}

	body, err := json.Marshal(&genReq{text, limit})
	if err != nil {
		return nil, err
	}

	rsp, err := svc.request(ctx, "POST", &url.URL{}, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	switch rsp.StatusCode {
	case http.StatusOK:
		return rsp.Body, nil

	default:
		io.Copy(ioutil.Discard, rsp.Body)
		rsp.Body.Close()
		return nil, errors.New(fmt.Sprintf("unexpected response: %d", rsp.StatusCode))
	}
}

//
// wordsvc: generates words to use when none are specified
//

// Make a client to wordsvc
func MakeWordSvc(scheme, addr string) WordSvc {
	return &wordSvc{Client{"word", scheme, addr, &http.Client{}}}
}

// Satisfy WordSvc with a Client
func (svc *wordSvc) Word(ctx contextHeaders) (string, error) {
	rsp, err := svc.request(ctx, "GET", &url.URL{}, nil)
	if err != nil {
		return "", err
	}

	switch rsp.StatusCode {
	case http.StatusOK:
		var word wordRsp
		if err := json.NewDecoder(rsp.Body).Decode(&word); err != nil {
			return "", err
		}
		io.Copy(ioutil.Discard, rsp.Body)
		rsp.Body.Close()
		return word.Word, nil

	default:
		io.Copy(ioutil.Discard, rsp.Body)
		rsp.Body.Close()
		return "", errors.New(fmt.Sprintf("unexpected response: %d", rsp.StatusCode))
	}
}

//
// main
//

func dieIf(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %s. Try --help for help.\n", err)
	os.Exit(-1)
}

func main() {
	addr := flag.String("srv", ":8080", "TCP address to listen on (in host:port form)")
	genAddr := flag.String("gen-addr", "localhost:8181", "Address of the gen service")
	wordAddr := flag.String("word-addr", "localhost:8282", "Address of the word service")
	flag.Parse()
	if flag.NArg() != 0 {
		dieIf(fmt.Errorf("expecting zero arguments but got %d", flag.NArg()))
	}

	server := &http.Server{
		Addr: *addr,
		Handler: &GobWeb{
			genSvc:  MakeGenSvc("http", *genAddr),
			wordSvc: MakeWordSvc("http", *wordAddr),
		},
	}
	dieIf(server.ListenAndServe())
}
