package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

var usageErrFmt = `reqz [--concurrency=1] [--method=GET] url
Error: %s

Try reqz --help for more info.
`

func exUsage(msg string) {
	fmt.Fprintf(os.Stderr, usageErrFmt, msg)
	os.Exit(64)
}

func main() {
	concurrency := flag.Uint("concurrency", 1, "Number of request threads")
	method := flag.String("method", "GET", "HTTP method")
	interval := flag.Duration("interval", 10*time.Second, "reporting interval")

	flag.Parse()
	if flag.NArg() != 1 {
		exUsage("expected a single comandline argument")
	}

	urlstr := flag.Arg(0)
	dstURL, err := url.Parse(urlstr)
	if err != nil {
		exUsage(fmt.Sprintf("invalid URL: '%s': %s", urlstr, err.Error()))
	}
	if dstURL.Scheme == "" {
		dstURL.Scheme = "http"
	}
	if dstURL.Host == "" {
		exUsage(fmt.Sprintf("invalid URL: '%s': no host", urlstr))
	}

	if *concurrency < 1 {
		exUsage("--concurrency must be at least 1")
	}

	count := uint64(0)
	size := uint64(0)
	timeout := time.After(*interval)

	client := &http.Client{}
	received := make(chan uint64)
	for i := uint(0); i != *concurrency; i++ {
		go func() {
			for {
				req, err := http.NewRequest(*method, dstURL.String(), nil)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					continue
				}

				rsp, err := client.Do(req)
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					continue
				}

				sz, _ := io.Copy(ioutil.Discard, rsp.Body)
				rsp.Body.Close()

				// inform the main thread how many bytes have been processed
				received <- uint64(sz)
			}
		}()
	}

	for {
		select {
		case t := <-timeout:
			// Periodically print stats about the health of the load generator
			fmt.Printf("%s: %d requests %d bytes / %s\n", t, count, size, interval)
			count = 0
			size = 0
			timeout = time.After(*interval)

		case bytes := <-received:
			// each time a request is processed, we increment the
			// request count and the number of bytes.
			count++
			size += bytes
		}
	}
}
