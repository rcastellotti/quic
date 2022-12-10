package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/lucas-clemente/quic-go/logging"
	"github.com/lucas-clemente/quic-go/qlog"
	"github.com/rcastellotti/quic-project/utils"
)

func setupHandler(www string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome, random ACN student :)\n"))
	})
	return mux
}

func main() {

	www := flag.String("www", "", "www data")
	enableQlog := flag.Bool("qlog", false, "output a qlog (in the same directory)")
	flag.Parse()

	handler := setupHandler(*www)
	quicConf := &quic.Config{}

	if *enableQlog {
		quicConf.Tracer = qlog.NewTracer(func(_ logging.Perspective, connID []byte) io.WriteCloser {
			filename := fmt.Sprintf("server_%x.qlog", connID)
			f, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Creating qlog file %s.\n", filename)
			return utils.NewBufferedWriteCloser(bufio.NewWriter(f), f)
		})
	}

	certFile := "../cert.pem"
	keyFile := "../priv.key"

	server := http3.Server{
		Handler:    handler,
		Addr:       "server:4433",
		QuicConfig: quicConf,
	}
	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		fmt.Println(err)
	}
}
