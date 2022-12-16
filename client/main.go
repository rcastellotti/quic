package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/lucas-clemente/quic-go/logging"
	"github.com/lucas-clemente/quic-go/qlog"
	"github.com/rcastellotti/quic-project/utils"
)

func main() {
	SSLKEYLOGFILE := os.Getenv("SSLKEYLOGFILE")
	QLOGDIR := os.Getenv("QLOGDIR")

	insecure := flag.Bool("insecure", false, "skip certificate verification")
	_, enableQlog := os.LookupEnv("QLOGDIR")

	flag.Parse()
	urls := flag.Args()
	// i think i can make this better
	var keyLog io.Writer
	f, err := os.Create(SSLKEYLOGFILE)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	keyLog = f

	pool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatal(err)
	}
	utils.AddRootCA(pool)

	var qconf quic.Config
	if enableQlog {
		qconf.Tracer = qlog.NewTracer(func(_ logging.Perspective, connID []byte) io.WriteCloser {
			filename := fmt.Sprintf(QLOGDIR+"/client_%x.qlog", connID)
			f, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Creating qlog file %s.\n", filename)
			return utils.NewBufferedWriteCloser(bufio.NewWriter(f), f)
		})
	}
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: *insecure,
			KeyLogWriter:       keyLog,
		},
		QuicConfig: &qconf,
	}
	defer roundTripper.Close()
	hclient := &http.Client{
		Transport: roundTripper,
	}

	rsp, err := hclient.Get(urls[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Got response for %s: %#v", urls[0], rsp)

	body := &bytes.Buffer{}
	_, err = io.Copy(body, rsp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response Body:")
	fmt.Printf("%s", body.Bytes())
}
