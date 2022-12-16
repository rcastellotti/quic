certs: ./certs.sh
client SSLKEYLOGFILE="key.log" QLOGDIR="." go run main.go -qlog https://server:4433
server: CERTS=$(dirname $PWD/) SERVERNAME="server" QLOGDIR="." WWW="/home/rc/github.com/rcastellotti/" PORT=4433 go run main.go 


tried to run this commands

SSLKEYLOGFILE="key.log" QLOGDIR="." go run main.go  https://server:4433/priv.key
SSLKEYLOGFILE="key.log" QLOGDIR="." go run main.go  https://localhost:4433/priv.key
SSLKEYLOGFILE="key.log" QLOGDIR="." go run main.go  https://127.0.0.1:4433/priv.key

