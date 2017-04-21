package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	pb "github.com/linkerd/linkerd-examples/gob/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type (
	GobWeb struct {
		genSvc  pb.GenSvcClient
		wordSvc pb.WordSvcClient
	}

	helpCtx struct {
		Host string
	}
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
	ctx := getContext(req)

	switch req.URL.Path {
	case "/", "/help":
		switch req.Method {
		case "GET":
			rspw.Header().Set("content-type", "text/plain")
			if err = HelpTemplate.Execute(rspw, &helpCtx{req.Host}); err != nil {
				fmt.Fprintf(os.Stderr, "template error: %s\n", err.Error())
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
				rsp, err := gob.wordSvc.GetWord(ctx, &pb.WordRequest{})
				if err != nil {
					fmt.Println("wordsvc error: " + err.Error())
					rspw.WriteHeader(http.StatusInternalServerError)
					return
				}
				text = rsp.Word
			}
			if text == "" {
				fmt.Println("could not load text")
				rspw.WriteHeader(http.StatusInternalServerError)
				return
			}

			limit := uint(0)
			limitstr := params.Get("limit")
			if limitstr != "" {
				limit32, err := strconv.ParseUint(limitstr, 10, 32)
				if err != nil {
					rspw.WriteHeader(http.StatusBadRequest)
					return
				}
				limit = uint(limit32)
			}

			stream, err := gob.genSvc.Gen(ctx, &pb.GenRequest{text, int32(limit)})
			if err != nil {
				fmt.Println("error generating: " + err.Error())
				rspw.WriteHeader(http.StatusInternalServerError)
				return
			}

			rspw.Header().Set("content-type", "text/plain")
			latency_param := params.Get("latency")
			if latency_param != "" {
				latency, err := time.ParseDuration(latency_param)
				if err != nil {
					rspw.WriteHeader(http.StatusBadRequest)
					return
				}
				time.Sleep(latency)
			}
			rspw.WriteHeader(http.StatusOK)

			streaming := true
			for streaming {
				rsp, err := stream.Recv()
				if err == io.EOF {
					streaming = false
				} else if err != nil {
					fmt.Println("streaming error: " + err.Error())
					streaming = false
				} else {
					if _, err := rspw.Write([]byte(rsp.Text)); err != nil {
						fmt.Println("write error: " + err.Error())
						streaming = false
						// XXX Close stream
					}
				}
			}
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

// extract router-specific headers to support dynamic routing & tracing
func getContext(req *http.Request) context.Context {
	headers := make(map[string]string)
	for k, values := range req.Header {
		prefixed := func(s string) bool { return strings.HasPrefix(k, s) }
		if prefixed("L5d-") || prefixed("Dtab-") || prefixed("X-Dtab-") {
			if len(values) > 0 {
				headers[k] = values[0]
			}
		}
	}
	md := metadata.New(headers)
	ctx := metadata.NewContext(context.Background(), md)
	return ctx
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
	caCertFile := flag.String("cacert", "", "path to PEM-formatted CA certificate")
	addr := flag.String("srv", ":8080", "TCP address to listen on (in host:port form)")
	genAddr := flag.String("gen-addr", "localhost:8181", "Address of the gen service")
	genName := flag.String("gen-name", "", "Common name of gen service")
	wordAddr := flag.String("word-addr", "localhost:8282", "Address of the word service")
	wordName := flag.String("word-name", "", "Common name of word service")
	flag.Parse()
	if flag.NArg() != 0 {
		dieIf(fmt.Errorf("expecting zero arguments but got %d", flag.NArg()))
	}

	var genCreds grpc.DialOption
	if *caCertFile == "" {
		genCreds = grpc.WithInsecure()
	} else {
		creds, err := credentials.NewClientTLSFromFile(*caCertFile, *genName)
		dieIf(err)
		genCreds = grpc.WithTransportCredentials(creds)
	}
	genConn, err := grpc.Dial(*genAddr, genCreds)
	dieIf(err)
	defer genConn.Close()
	genClient := pb.NewGenSvcClient(genConn)

	var wordCreds grpc.DialOption
	if *caCertFile == "" {
		wordCreds = grpc.WithInsecure()
	} else {
		creds, err := credentials.NewClientTLSFromFile(*caCertFile, *wordName)
		dieIf(err)
		wordCreds = grpc.WithTransportCredentials(creds)
	}
	wordConn, err := grpc.Dial(*wordAddr, wordCreds)
	dieIf(err)
	defer wordConn.Close()
	wordClient := pb.NewWordSvcClient(wordConn)

	server := &http.Server{
		Addr: *addr,
		Handler: &GobWeb{
			genSvc:  genClient,
			wordSvc: wordClient,
		},
	}
	dieIf(server.ListenAndServe())
}
