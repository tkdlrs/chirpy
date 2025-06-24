package chirpy

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"
)

type Server struct {
	Addr string

	Handler Handler // handler to invoke, http.DefaultServeMux if nil

	DisableGeneralOptionsHandler bool

	TLSConfig *tls.Config

	ReadTimeout time.Duration

	ReadHeaderTimeout time.Duration

	WriteTimeout time.Duration

	IdleTimeout time.Duration

	MaxHeaderBytes int

	TLSNextProto map[string]func(*Server, *tls.Conn, Handler)

	ConnState func(net.Conn, ConnState)

	ErrorLog *log.Logger

	BaseContext func(net.Listener) context.Context

	ConnContext func(ctx context.Context, c net.Conn) context.Context

	HTTP2 *HTTP2Config

	Protocols *Protocols
}

func main() {
	mux := http.NewServeMux
}
