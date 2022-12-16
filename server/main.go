package main

import (
	"bufio"
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

func setupHandler(www string) http.Handler {
	fmt.Print(http.Dir(www))
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(www)))

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome, random ACN student :)\n"))
	})

	return mux
}

func main() {
	SERVERNAME := os.Getenv("SERVERNAME")
	PORT := os.Getenv("PORT")
	CERTS := os.Getenv("CERTS")
	QLOGDIR := os.Getenv("QLOGDIR")
	WWW := os.Getenv("WWW")

	serverAndPort := fmt.Sprintf(SERVERNAME + ":" + PORT)
	_, enableQlog := os.LookupEnv("QLOGDIR")
	flag.Parse()

	handler := setupHandler(WWW)
	quicConf := &quic.Config{}

	if enableQlog {
		quicConf.Tracer = qlog.NewTracer(func(_ logging.Perspective, connID []byte) io.WriteCloser {
			filename := fmt.Sprintf(QLOGDIR+"/server_%x.qlog", connID)
			f, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Creating qlog file %s.\n", filename)
			return utils.NewBufferedWriteCloser(bufio.NewWriter(f), f)
		})
	}

	certFile := fmt.Sprintf(CERTS + "/cert.pem")
	keyFile := fmt.Sprintf(CERTS + "/priv.key")
	server := http3.Server{
		Handler:    handler,
		Addr:       serverAndPort,
		QuicConfig: quicConf,
	}
	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		fmt.Println(err)
	}
}
