package proxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"

	"github.com/things-go/go-socks5"
)

var listener net.Listener

type Config struct {
	SocksListenAddr string
	HttpListenAddr  string
	Mappings        map[string]int // key: domain, value: local port
}

func Start(cfg *Config) error {
	if cfg.SocksListenAddr == "" {
		return nil
	}
	l, err := net.Listen("tcp", cfg.SocksListenAddr)
	if err != nil {
		return fmt.Errorf("cannot start TCP listener on %s: %w", cfg.SocksListenAddr, err)
	}

	localNames := make([]string, len(cfg.Mappings))
	domainIdx := 0
	for domain := range cfg.Mappings {
		localNames[domainIdx] = domain
		domainIdx++
	}

	listener = l // using own listener to be able to run async
	socksServer := socks5.NewServer(socks5.WithResolver(&resolver{
		localNames: localNames,
	}))
	log.Printf("Started proxy server on %s\n", cfg.SocksListenAddr)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Error accepting new proxy connection: %s\n", err.Error())
				continue
			}
			go func() {
				defer conn.Close()
				if err = socksServer.ServeConn(conn); err != nil {
					log.Printf("error socks proxying a connection: %s\n", err.Error())
				}
			}()
		}
	}()

	if cfg.HttpListenAddr == "" {
		return nil
	}

	localDirector := createLocalReverseDirector(cfg.Mappings)
	reverseProxy := &httputil.ReverseProxy{
		Director: localDirector,
	}
	go func() {
		log.Printf("starting http reverse proxy at %s\n", cfg.HttpListenAddr)
		if err := http.ListenAndServe(cfg.HttpListenAddr, reverseProxy); err != nil {
			log.Printf("http reverse proxy failed: %s\n", err.Error())
		}
	}()

	return nil
}

func Stop() {
}
